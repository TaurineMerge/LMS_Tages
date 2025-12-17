"""JWT security and token validation."""

import logging
from datetime import datetime
from typing import Optional

import httpx
from fastapi import Depends, HTTPException, Request, status
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from jose import jwt
from jose.exceptions import ExpiredSignatureError, JWTError
from pydantic import BaseModel, Field

from app.config import get_settings
from app.core.jwt import JwtService
from app.telemetry import traced

logger = logging.getLogger(__name__)
settings = get_settings()

# Security scheme for FastAPI
# allow optional credential extraction so we can fallback to HttpOnly cookie
security = HTTPBearer(auto_error=False)


class TokenPayload(BaseModel):
    """JWT token payload from Keycloak."""

    sub: str = Field(..., description="User ID (subject)")
    email: Optional[str] = None
    given_name: Optional[str] = None  # ← добавлено
    family_name: Optional[str] = None
    preferred_username: Optional[str] = None
    exp: int = Field(..., description="Token expiration time")
    iat: Optional[int] = None
    roles: list[str] = Field(default_factory=list)

    class Config:
        json_schema_extra = {
            "example": {
                "sub": "user-id-uuid",
                "email": "user@example.com",
                "given_name": "User",
                "family_name": "Name",
                "preferred_username": "username",
                "exp": 1765284914,
                "iat": 1765284614,
            }
        }


# replace internal jwks/cache with shared service
_jwt_service = JwtService(keycloak_url=settings.KEYCLOAK_SERVER_URL, realm=settings.KEYCLOAK_REALM)


class JWTValidator:
    """Validate JWT tokens using Keycloak JWKS."""

    def __init__(self, keycloak_url: str, issuer_url: str, realm: str, client_id: str):
        self.client_id = client_id
        self.issuer = f"{issuer_url.rstrip('/')}/realms/{realm}"  # ← внешний URL для issuer
        self.keycloak_url = keycloak_url  # ← внутренний URL для получения ключей

    @traced("jwt_validator.validate_token", record_args=True, record_result=True)
    async def validate_token(self, token: str) -> TokenPayload:
        """
        Validate JWT token and extract payload.

        Args:
            token: JWT token string

        Returns:
            TokenPayload with user data

        Raises:
            HTTPException: If token is invalid or expired
        """
        try:
            payload = await _jwt_service.decode(
                token, audience="account", issuer=self.issuer, keycloak_url=self.keycloak_url
            )
            # extract roles as before
            roles = []
            if isinstance(payload, dict):
                realm_access = payload.get("realm_access")
                if isinstance(realm_access, dict):
                    roles = realm_access.get("roles", []) or []
                # fallback: some setups may put roles directly
                if not roles:
                    roles = payload.get("roles", []) or []

                # ensure serializable list
                payload["roles"] = list(roles)
            # Преобразуем в типизированный объект
            token_payload = TokenPayload(**payload)

            logger.info(f"Token validated successfully for user {token_payload.sub}")
            return token_payload

        except ExpiredSignatureError:
            logger.warning("JWT token has expired")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Token has expired",
                headers={"WWW-Authenticate": "Bearer"},
            )
        except JWTError as e:
            logger.warning(f"JWT validation failed: {e}")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid token", headers={"WWW-Authenticate": "Bearer"}
            )
        except Exception as e:
            logger.error(f"Unexpected error during token validation: {e}")
            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Token validation error")


# Создаём singleton validator
_jwt_validator = JWTValidator(
    keycloak_url=settings.KEYCLOAK_SERVER_URL,
    issuer_url=settings.KEYCLOAK_PUBLIC_URL,
    realm=settings.KEYCLOAK_REALM,
    client_id=settings.KEYCLOAK_CLIENT_ID,
)


@traced("security.get_current_user", record_result=False)
async def get_current_user(
    request: Request,
    credentials: Optional[HTTPAuthorizationCredentials] = Depends(security),
) -> TokenPayload:
    """
    Get current authenticated user from JWT token.

    Этот dependency используется в защищённых роутах:

    ```python
    @router.get("/me")
    async def get_me(current_user: TokenPayload = Depends(get_current_user)):
        return {"user": current_user}
    ```
    """
    token = None

    # Prefer Authorization header if present
    if credentials and credentials.credentials:
        token = credentials.credentials

    # Fallback to HttpOnly cookie set by the callback endpoint
    if not token:
        cookie_token = request.cookies.get("access_token")
        if cookie_token:
            token = cookie_token

    if not token:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail="Not authenticated", headers={"WWW-Authenticate": "Bearer"}
        )

    return await _jwt_validator.validate_token(token)


@traced("security.get_current_user_optional", record_result=False)
async def get_current_user_optional(
    request: Request,
    credentials: Optional[HTTPAuthorizationCredentials] = Depends(security),
) -> Optional[TokenPayload]:
    """
    Get current user if authenticated, otherwise return None.

    Используется для опциональной аутентификации.
    """
    token = None

    if credentials and credentials.credentials:
        token = credentials.credentials

    if not token:
        token = request.cookies.get("access_token")

    if not token:
        return None

    try:
        return await _jwt_validator.validate_token(token)
    except HTTPException:
        # If token invalid, return None
        return None


def require_roles(*roles: str):
    """
    Create a dependency that requires specific roles.

    ```python
    @router.get("/admin")
    async def admin_only(
        current_user: TokenPayload = Depends(get_current_user),
        _: None = Depends(require_roles("admin"))
    ):
        return {"message": "Admin only"}
    ```
    """

    async def check_roles(current_user: TokenPayload = Depends(get_current_user)) -> None:
        # Ensure the user is authenticated
        if not current_user:
            raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Not authenticated")

        # Check that at least one of the required roles is present in token
        if roles:
            user_roles = set(getattr(current_user, "roles", []) or [])
            required = set(roles)
            if user_roles.isdisjoint(required):
                raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Insufficient permissions")

    return check_roles


# Re-export for backward compatibility
__all__ = [
    "get_current_user",
    "get_current_user_optional",
    "require_roles",
    "TokenPayload",
    "JWTValidator",
    "security",
]
