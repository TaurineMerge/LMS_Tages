"""Authentication service using Keycloak."""
import logging
from typing import Dict, Any
from fastapi import HTTPException, status
from fastapi.concurrency import run_in_threadpool

from app.config import get_settings
from app.services.keycloak import keycloak_service
from app.telemetry import traced

logger = logging.getLogger(__name__)
settings = get_settings()

class AuthService:
    """High-level authentication logic."""

    @traced("auth.get_login_url", record_result=True)
    def get_login_url(self) -> str:
        """Generate Keycloak login URL."""
        redirect_uri = settings.KEYCLOAK_REDIRECT_URI
        return keycloak_service.get_auth_url(redirect_uri)

    @traced("auth.exchange_code", record_args=True, record_result=True)
    async def exchange_code_for_token(self, code: str) -> Dict[str, Any]:
        """Exchange authorization code for access token."""
        redirect_uri = settings.KEYCLOAK_REDIRECT_URI

        try:
            return await run_in_threadpool(
                keycloak_service.get_token,
                code,
                redirect_uri
            )
        except Exception as e:
            logger.error(f"Token exchange failed: {e}")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Failed to exchange code for token"
            )

    @traced("auth.refresh", record_args=True, record_result=True)
    async def refresh_token(self, refresh_token: str) -> Dict[str, Any]:
        """Refresh access token using refresh token."""
        try:
            return await run_in_threadpool(
                keycloak_service.refresh_token,
                refresh_token
            )
        except Exception as e:
            logger.error(f"Token refresh failed: {e}")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid refresh token"
            )

    @traced("auth.logout")
    async def logout(self, refresh_token: str) -> None:
        """Logout user."""
        try:
            await run_in_threadpool(keycloak_service.logout, refresh_token)
        except Exception as e:
            logger.warning(f"Logout failed: {e}")
            # Logout is best-effort, don't raise

    @traced("auth.register_user", record_args=True, record_result=True)
    async def register_user(self, username: str, email: str, password: str, first_name: str, last_name: str) -> Dict[str, str]:
        """Registers a user in Keycloak."""
        user_data = {
            "username": username,
            "email": email,
            "firstName": first_name,
            "lastName": last_name,
            "enabled": True,
            "credentials": [{"value": password, "type": "password", "temporary": False}],
            "emailVerified": settings.KEYCLOAK_USER_EMAIL_VERIFIED_DEFAULT
        }

        try:
            # Создаем пользователя в Keycloak
            user_id = await run_in_threadpool(
                keycloak_service.create_user,
                user_data
            )
            
            logger.info(f"User registered successfully: {username}")

            return {
                "user_id": user_id,
                "username": username,
                "email": email
            }
        except Exception as e:
            # Обработка ошибок Keycloak
            error_msg = str(e)
            if "409" in error_msg or "already exists" in error_msg.lower():
                raise HTTPException(status_code=409, detail="User already exists")
            
            logger.error(f"Registration failed: {e}")
            raise HTTPException(status_code=500, detail="Registration failed")

auth_service = AuthService()
