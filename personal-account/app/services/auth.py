"""Authentication service using Keycloak."""

import logging
import uuid
from typing import Any, Dict
from urllib.parse import parse_qs, urlencode, urlparse

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

    # @traced("auth.get_register_url", record_args=True, record_result=True)
    # def get_register_url(self, redirect_uri: str) -> str:
    #     """Generate Keycloak registration URL using existing keycloak_service."""
    #     try:
    #         # Получаем базовый auth URL через существующий keycloak_service
    #         auth_url = keycloak_service.get_auth_url(redirect_uri)

    #         # Парсим URL для безопасного добавления параметров
    #         parsed_url = urlparse(auth_url)
    #         query_params = parse_qs(parsed_url.query)

    #         # Добавляем параметр для перехода на регистрацию
    #         query_params["kc_action"] = ["register"]

    #         # Генерируем уникальный state с помощью uuid
    #         query_params["state"] = [str(uuid.uuid4())]

    #         # Собираем URL обратно
    #         new_query = urlencode(query_params, doseq=True)
    #         registration_url = parsed_url._replace(query=new_query).geturl()

    #         logger.debug(f"Keycloak registration URL generated: {registration_url}")
    #         return registration_url

    #     except Exception as e:
    #         logger.error(f"Failed to generate registration URL: {str(e)}")
    #         raise

    @traced("auth.logout", record_args=True, record_result=True)
    async def logout(self, refresh_token: str) -> None:
        """Logout user."""
        try:
            await run_in_threadpool(keycloak_service.logout, refresh_token)
        except Exception as e:
            logger.warning(f"Logout failed: {e}")
            # Logout is best-effort, don't raise


auth_service = AuthService()
