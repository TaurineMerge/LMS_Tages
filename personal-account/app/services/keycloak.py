import logging
from typing import Any, Dict, Optional

from keycloak.exceptions import KeycloakAuthenticationError, KeycloakError

from app.config import get_settings
from app.telemetry import traced
from keycloak import KeycloakAdmin, KeycloakOpenID

logger = logging.getLogger(__name__)
settings = get_settings()


class KeycloakService:
    """Service for direct interactions with Keycloak API."""

    def __init__(self):
        # 1. Настройка OpenID (для логина пользователей)
        self.openid = KeycloakOpenID(
            server_url=settings.KEYCLOAK_SERVER_URL,
            client_id=settings.KEYCLOAK_CLIENT_ID,
            realm_name=settings.KEYCLOAK_REALM,
            client_secret_key=settings.KEYCLOAK_CLIENT_SECRET,
            verify=True,
        )

        # 2. Настройка Admin (для управления пользователями)
        # Важно: инициализируем лениво или сразу, но с обработкой ошибок
        self._admin_client = None

    @property
    def admin(self) -> KeycloakAdmin:
        """Lazy initialization of Keycloak Admin client."""
        if not self._admin_client:
            try:
                self._admin_client = KeycloakAdmin(
                    server_url=settings.KEYCLOAK_SERVER_URL,
                    username=settings.KEYCLOAK_ADMIN_USERNAME,
                    password=settings.KEYCLOAK_ADMIN_PASSWORD,
                    realm_name=settings.KEYCLOAK_REALM,
                    user_realm_name="master",  # Админ обычно в master реалме
                    verify=True,
                )
            except Exception as e:
                logger.error(f"Failed to initialize Keycloak Admin: {e}")
                raise
        return self._admin_client

    @traced("keycloak.get_auth_url", record_args=True, record_result=True)
    def get_auth_url(self, redirect_uri: str) -> str:
        return self.openid.auth_url(redirect_uri=redirect_uri)

    @traced("keycloak.get_token", record_args=True, record_result=True)
    def get_token(self, code: str, redirect_uri: str) -> Dict[str, Any]:
        return self.openid.token(grant_type="authorization_code", code=code, redirect_uri=redirect_uri)

    @traced("keycloak.refresh_token", record_args=True, record_result=True)
    def refresh_token(self, refresh_token: str) -> Dict[str, Any]:
        return self.openid.refresh_token(refresh_token)

    @traced("keycloak.logout", record_args=True)
    def logout(self, refresh_token: str) -> None:
        self.openid.logout(refresh_token)

    @traced("keycloak.create_user", record_args=True, record_result=True)
    def create_user(self, user_data: Dict[str, Any]) -> str:
        """Creates a user and returns their ID."""
        return self.admin.create_user(user_data)

    @traced("keycloak.get_user_id", record_args=True, record_result=True)
    def get_user_id(self, username: str) -> Optional[str]:
        return self.admin.get_user_id(username)


# Singleton instance
keycloak_service = KeycloakService()
