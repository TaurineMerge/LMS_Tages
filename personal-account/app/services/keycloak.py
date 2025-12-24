import logging
from typing import Any, Dict, Optional

from keycloak import KeycloakAdmin, KeycloakOpenID
from keycloak.exceptions import KeycloakError

from app.config import get_settings
from app.telemetry import traced

logger = logging.getLogger(__name__)
settings = get_settings()


class KeycloakService:
    """Service for direct interactions with Keycloak API."""

    def __init__(self):
        self.verify_ssl = False
        # Логируем настройки
        logger.info(f"KeycloakService initialized with SSL verification: {self.verify_ssl}")
        logger.info(f"Keycloak server URL: {settings.KEYCLOAK_SERVER_URL}")

        # 1. Настройка OpenID (для логина пользователей)
        self.openid = KeycloakOpenID(
            server_url=settings.KEYCLOAK_SERVER_URL,
            client_id=settings.KEYCLOAK_CLIENT_ID,
            realm_name=settings.KEYCLOAK_REALM,
            client_secret_key=settings.KEYCLOAK_CLIENT_SECRET,
            verify=self.verify_ssl,
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
                    verify=False,
                    custom_headers={"User-Agent": "python-keycloak-dev"},
                )
            except Exception as e:
                logger.error(f"Failed to initialize Keycloak Admin: {e}")
                raise
        return self._admin_client

    @traced("keycloak.get_user_data_by_token")
    def get_user_data_by_token(self, access_token: str) -> Dict[str, Any]:
        try:
            # Получаем userinfo напрямую через OpenID Connect
            info = self.openid.userinfo(token=access_token)
            return {
                "id": info.get("sub"),  # ← sub — это user ID в Keycloak
                "name": info.get("given_name", ""),
                "surname": info.get("family_name", ""),
                "email": info.get("email", ""),
                "username": info.get("preferred_username", ""),
            }
        except KeycloakError as e:
            logger.error(f"Failed to fetch userinfo: {e}")
            raise

    @traced("keycloak.get_auth_url", record_args=True, record_result=True)
    def get_auth_url(self, redirect_uri: str) -> str:
        try:
            url = self.openid.auth_url(
                redirect_uri=redirect_uri,
                scope="openid profile email",  # ← ОБЯЗАТЕЛЬНО
            )
            logger.info(f"Keycloak auth URL: {url}")
            return url
        except Exception as e:
            logger.error(f"Keycloak auth URL generation failed: {e}")
            raise

    @traced("keycloak.get_token", record_args=True, record_result=True)
    def get_token(self, code: str, redirect_uri: str) -> Dict[str, Any]:
        tokens = self.openid.token(
            grant_type="authorization_code",
            code=code,
            redirect_uri=redirect_uri,
            scope="openid profile email",  # ← ДОБАВЛЕНО
        )
        logger.debug(f"Token exchange successful. Scopes: {tokens.get('scope', 'N/A')}")
        return tokens

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
