import logging
from typing import Any, Optional

from keycloak import KeycloakAdmin, KeycloakOpenID
from keycloak.exceptions import KeycloakAuthenticationError, KeycloakError

from app.config import get_settings
from app.telemetry import traced

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
        try:
            url = self.openid.auth_url(redirect_uri=redirect_uri)
            print(f"Keycloak auth URL: {url}")  # Добавь для отладки
            return url
        except Exception as e:
            print(f"Keycloak error: {e}")  # Добавь
            raise

    # @traced("keycloak.get_registration_url", record_args=True, record_result=True)
    # def get_registration_url(self, redirect_uri: str) -> str:
    #     try:
    #         # Сначала получаем базовый URL аутентификации
    #         auth_url = self.openid.auth_url(redirect_uri=redirect_uri)

    #         # Заменяем 'auth' на 'registrations' в пути URL
    #         # Например: /realms/myrealm/protocol/openid-connect/auth
    #         # станет: /realms/myrealm/protocol/openid-connect/registrations
    #         registration_url = auth_url.replace("/auth?", "/registrations?")

    #         print(f"Keycloak registration URL: {registration_url}")
    #         return registration_url
    #     except Exception as e:
    #         print(f"Keycloak registration error: {e}")
    #         raise

    @traced("keycloak.get_user_data", record_args=True, record_result=True)
    def get_user_data(self, user_id: str) -> dict[str, Any]:
        try:
            user_info = self.admin.get_user(user_id)
            return {
                "id": user_info.get("id"),
                "name": user_info.get("firstName", ""),
                "surname": user_info.get("lastName", ""),
                "email": user_info.get("email", ""),
                "username": user_info.get("username", ""),
                # при необходимости добавьте остальные поля:
                # "phone": user_info.get("attributes", {}).get("phoneNumber", [None])[0],
                # "enabled": user_info.get("enabled", False),
                # "emailVerified": user_info.get("emailVerified", False),
            }
        except KeycloakError as e:
            logger.error(f"Keycloak API error during user fetch: {e}")
            raise
        except Exception as e:
            logger.error(f"Unexpected error during user fetch from Keycloak: {e}")
            raise

    @traced("keycloak.get_token", record_args=True, record_result=True)
    def get_token(self, code: str, redirect_uri: str) -> dict[str, Any]:
        return self.openid.token(grant_type="authorization_code", code=code, redirect_uri=redirect_uri)

    @traced("keycloak.refresh_token", record_args=True, record_result=True)
    def refresh_token(self, refresh_token: str) -> dict[str, Any]:
        return self.openid.refresh_token(refresh_token)

    @traced("keycloak.logout", record_args=True)
    def logout(self, refresh_token: str) -> None:
        self.openid.logout(refresh_token)

    @traced("keycloak.create_user", record_args=True, record_result=True)
    def create_user(self, user_data: dict[str, Any]) -> str:
        """Creates a user and returns their ID."""
        return self.admin.create_user(user_data)

    @traced("keycloak.get_user_id", record_args=True, record_result=True)
    def get_user_id(self, username: str) -> Optional[str]:
        return self.admin.get_user_id(username)

    @traced("keycloak.update_user_data", record_args=True, record_result=True)
    def update_user_data(self, user_id: str, data: dict[str, Any]):
        """Updates user data in Keycloak."""
        try:
            logger.info(f"🔄 STARTING update_user_data for user_id: {user_id}")
            logger.info(f"📥 Input data: {data}")

            # Получаем полный объект пользователя
            user_obj = self.admin.get_user(user_id)
            logger.info(f"📋 Current user object: {user_obj}")

            # Обновляем только указанные поля
            if "firstName" in data and data["firstName"] is not None:
                user_obj["firstName"] = data["firstName"]
            if "lastName" in data and data["lastName"] is not None:
                user_obj["lastName"] = data["lastName"]
            if "email" in data and data["email"] is not None:
                user_obj["email"] = data["email"]
            if "username" in data and data["username"] is not None:
                user_obj["username"] = data["username"]  # Пытаемся обновить через полный объект

            logger.info(f"📦 Updated user object: {user_obj}")

            # Обновляем весь объект
            result = self.admin.update_user(user_id=user_id, payload=user_obj)
            logger.info(f"✅ Keycloak API update result: {result}")

            # Проверяем результат
            updated_user = self.admin.get_user(user_id)
            logger.info(f"📋 User data AFTER full update: {updated_user}")

            # Проверяем, что изменения применились
            changes_applied = True
            if "firstName" in data and updated_user.get("firstName") != data["firstName"]:
                logger.error(
                    f"❌ Name was not updated! Expected: {data['firstName']}, Got: {updated_user.get('firstName')}"
                )
                changes_applied = False
            if "lastName" in data and updated_user.get("lastName") != data["lastName"]:
                logger.error(
                    f"❌ Surname was not updated! Expected: {data['lastName']}, Got: {updated_user.get('lastName')}"
                )
                changes_applied = False
            if "email" in data and updated_user.get("email") != data["email"]:
                logger.error(f"❌ Email was not updated! Expected: {data['email']}, Got: {updated_user.get('email')}")
                changes_applied = False
            if "username" in data and updated_user.get("username") != data["username"]:
                logger.error(
                    f"❌ Username was not updated! Expected: {data['username']}, Got: {updated_user.get('username')}"
                )
                changes_applied = False

            if changes_applied:
                logger.info(f"🎉 User {user_id} updated successfully in Keycloak.")
            else:
                logger.error(f"💥 Some fields were not updated for user {user_id}")

        except Exception as e:
            logger.error(f"💥 ERROR in update_user_data: {e}", exc_info=True)
            raise


# Singleton instance
keycloak_service = KeycloakService()
