"""Authentication service using Keycloak."""
import logging
from keycloak import KeycloakOpenID, KeycloakAdmin
from keycloak.exceptions import KeycloakAuthenticationError, KeycloakGetError
from fastapi import HTTPException, status, Depends
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials

from app.config import get_settings
from app.telemetry import traced

logger = logging.getLogger(__name__)
settings = get_settings()
security = HTTPBearer()

class auth_service:
    """Service for handling Keycloak authentication."""

    def __init__(self):
        self.keycloak_openid = KeycloakOpenID(
            server_url=settings.KEYCLOAK_SERVER_URL,
            client_id=settings.KEYCLOAK_CLIENT_ID,
            realm_name=settings.KEYCLOAK_REALM,
            client_secret_key=settings.KEYCLOAK_CLIENT_SECRET,
        )
        self._keycloak_admin = None

    @property
    @traced("auth_service.keycloak_admin")
    def keycloak_admin(self) -> KeycloakAdmin:
        """Get or create KeycloakAdmin instance."""
        if self._keycloak_admin is None:
            # Admin authenticates in master realm, then switches to target realm
            self._keycloak_admin = KeycloakAdmin(
                server_url=settings.KEYCLOAK_SERVER_URL,
                username=settings.KEYCLOAK_ADMIN_USERNAME,
                password=settings.KEYCLOAK_ADMIN_PASSWORD,
                realm_name=settings.KEYCLOAK_ADMIN_REALM,  # Authenticate in master realm
                verify=True
            )
            # Switch to target realm for user management
            self._keycloak_admin.realm_name = settings.KEYCLOAK_REALM
        return self._keycloak_admin

    @traced("auth_service.get_login_url")
    def get_login_url(self) -> str:
        """Generate Keycloak login URL."""
        url = self.keycloak_openid.auth_url(
            redirect_uri=settings.KEYCLOAK_REDIRECT_URI,
            scope=settings.KEYCLOAK_DEFAULT_SCOPE,
            state="some_random_state"  # In production, use a secure random state
        )
        # Log generated URL for troubleshooting client lookup issues
        try:
            logger.debug("Generated Keycloak auth URL (raw): %s", url)
        except Exception:
            pass
        # Replace internal URL with public URL for browser redirection
        if settings.KEYCLOAK_SERVER_URL != settings.KEYCLOAK_PUBLIC_URL:
            return url.replace(settings.KEYCLOAK_SERVER_URL, settings.KEYCLOAK_PUBLIC_URL)
        return url

    @traced("auth_service.get_register_url")
    def get_register_url(self) -> str:
        """Generate Keycloak registration URL."""
        # For Keycloak 21+, use kc_action=register parameter with auth URL
        # The /registrations endpoint is deprecated
        base_url = settings.KEYCLOAK_PUBLIC_URL
        realm = settings.KEYCLOAK_REALM
        client_id = settings.KEYCLOAK_CLIENT_ID
        redirect_uri = settings.KEYCLOAK_REDIRECT_URI
        
        # Use standard auth endpoint with kc_action=register
        register_url = (
            f"{base_url}/realms/{realm}/protocol/openid-connect/auth"
            f"?client_id={client_id}"
            f"&redirect_uri={redirect_uri}"
            f"&response_type=code"
            f"&scope={settings.KEYCLOAK_DEFAULT_SCOPE}"
            f"&kc_action=register"
        )
        return register_url

    @traced("auth_service.register_user")
    async def register_user(
        self,
        username: str,
        email: str,
        password: str,
        first_name: str | None = None,
        last_name: str | None = None
    ) -> dict:
        """Register a new user in Keycloak."""
        try:
            user_data = {
                "username": username,
                "email": email,
                "enabled": True,
                "emailVerified": settings.KEYCLOAK_USER_EMAIL_VERIFIED_DEFAULT,  # Set to False if you want email verification
                "credentials": [{
                    "type": "password",
                    "value": password,
                    "temporary": False
                }]
            }
            
            if first_name:
                user_data["firstName"] = first_name
            if last_name:
                user_data["lastName"] = last_name
            
            # Create user
            user_id = self.keycloak_admin.create_user(user_data)
            
            logger.info(f"User registered successfully: {username} (ID: {user_id})")
            
            return {
                "user_id": user_id,
                "username": username,
                "email": email
            }
            
        except KeycloakGetError as e:
            error_message = str(e)
            if "User exists" in error_message or "409" in error_message:
                raise HTTPException(
                    status_code=status.HTTP_409_CONFLICT,
                    detail="User with this username or email already exists"
                )
            logger.error(f"Failed to register user: {error_message}")
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=f"Failed to register user: {error_message}"
            )
        except Exception as e:
            logger.error(f"Unexpected error during registration: {str(e)}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail=f"Registration failed: {str(e)}"
            )

    @traced("auth_service.exchange_code_for_token")
    async def exchange_code_for_token(self, code: str) -> dict:
        """Exchange authorization code for access token."""
        try:
            logger.info(f"Attempting to exchange authorization code (length: {len(code)})")
            logger.debug(f"Redirect URI: {settings.KEYCLOAK_REDIRECT_URI}")
            logger.debug(f"Keycloak server URL: {settings.KEYCLOAK_SERVER_URL}")
            
            token = self.keycloak_openid.token(
                grant_type="authorization_code",
                code=code,
                redirect_uri=settings.KEYCLOAK_REDIRECT_URI
            )
            logger.info("Code exchange successful, token obtained")
            return token
        except KeycloakAuthenticationError as e:
            logger.error(f"Keycloak auth error during code exchange: {str(e)}")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail=f"Failed to exchange code: {str(e)}"
            )
        except Exception as e:
            logger.error(f"Unexpected error during code exchange: {str(e)}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail=f"Failed to exchange code: {str(e)}"
            )

    @traced("auth_service.get_token_by_password")
    async def get_token_by_password(self, username: str, password: str) -> dict:
        """Get access token using username and password (OAuth2 password flow)."""
        try:
            token = self.keycloak_openid.token(
                grant_type="password",
                username=username,
                password=password
            )
            return token
        except KeycloakAuthenticationError as e:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid user credentials"
            )
        except Exception as e:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail=f"Failed to obtain token: {str(e)}"
            )

    @traced("auth_service.verify_token")
    async def verify_token(self, token: str) -> dict:
        """Verify and decode JWT token."""
        try:
            # Verify signature and expiration
            # options={"verify_signature": True, "verify_aud": True, "exp": True}
            user_info = self.keycloak_openid.userinfo(token)
            return user_info
        except Exception as e:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid authentication credentials"
            )

    @traced("auth_service.refresh_token")
    async def refresh_token(self, refresh_token: str) -> dict:
        """Refresh access token using refresh token."""
        try:
            logger.debug(f"Attempting to refresh token (token length: {len(refresh_token)})")
            # Use internal server URL for refresh to avoid issuer mismatch
            # Tokens issued by public URL have public issuer, but Keycloak
            # may expect internal issuer for refresh grant
            keycloak_for_refresh = KeycloakOpenID(
                server_url=settings.KEYCLOAK_SERVER_URL,
                client_id=settings.KEYCLOAK_CLIENT_ID,
                realm_name=settings.KEYCLOAK_REALM,
                client_secret_key=settings.KEYCLOAK_CLIENT_SECRET,
            )
            logger.debug(f"Keycloak refresh client configured: server={settings.KEYCLOAK_SERVER_URL}, realm={settings.KEYCLOAK_REALM}")
            token = keycloak_for_refresh.refresh_token(refresh_token)
            logger.info("Token refresh successful")
            return token
        except KeycloakAuthenticationError as e:
            logger.warning(f"Token refresh failed (auth error): {str(e)}")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid or expired refresh token"
            )
        except Exception as e:
            logger.error(f"Token refresh failed: {str(e)}")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail=f"Failed to refresh token: {str(e)}"
            )

    @traced("auth_service.logout")
    async def logout(self, refresh_token: str) -> None:
        """Logout user by invalidating the refresh token."""
        try:
            self.keycloak_openid.logout(refresh_token)
        except Exception as e:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=f"Failed to logout: {str(e)}"
            )


# Singleton instance
auth_service = auth_service()

async def get_current_user(credentials: HTTPAuthorizationCredentials = Depends(security)) -> dict:
    """Validate access token and return user info."""
    token = credentials.credentials
    return await auth_service.verify_token(token)
