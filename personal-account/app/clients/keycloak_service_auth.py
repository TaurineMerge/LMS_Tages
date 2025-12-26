"""Keycloak service-to-service authentication.

This module provides JWT token acquisition for service-to-service
communication with internal APIs protected by Keycloak authentication.

Uses either:
1. Client Credentials Grant (service account) - preferred for S2S
2. Resource Owner Password Grant (test user credentials) - fallback

Example
-------
```py
from app.clients.keycloak_service_auth import KeycloakServiceAuth

auth = KeycloakServiceAuth()
token = await auth.get_access_token()
headers = {"Authorization": f"Bearer {token}"}
```
"""

from __future__ import annotations

import logging
import time
from dataclasses import dataclass
from typing import Any

import httpx

from app.config import get_settings

logger = logging.getLogger(__name__)


@dataclass
class TokenCache:
    """Simple in-memory token cache with expiry."""

    access_token: str = ""
    expires_at: float = 0.0

    def is_valid(self) -> bool:
        """Check if cached token is still valid (with 30s buffer)."""
        return self.access_token and time.time() < (self.expires_at - 30)


class KeycloakServiceAuth:
    """Keycloak authentication client for service-to-service communication.

    This client obtains JWT tokens from Keycloak to authenticate
    requests to other services (like Testing service) that require
    Keycloak JWT authentication.

    Configuration is read from app settings (via environment variables):
    - KEYCLOAK_SERVER_URL: Keycloak base URL
    - KEYCLOAK_REALM: Target realm (default: student)
    - KEYCLOAK_SERVICE_CLIENT_ID: Client ID for S2S auth
    - KEYCLOAK_SERVICE_CLIENT_SECRET: Client secret

    For development, can also use username/password credentials.
    """

    def __init__(self) -> None:
        self.settings = get_settings()
        self._cache = TokenCache()

        # Service account settings (primary method)
        self.server_url = self.settings.KEYCLOAK_SERVER_URL
        self.realm = getattr(self.settings, "KEYCLOAK_SERVICE_REALM", self.settings.KEYCLOAK_REALM)
        self.client_id = getattr(self.settings, "KEYCLOAK_SERVICE_CLIENT_ID", "student-client")
        self.client_secret = getattr(self.settings, "KEYCLOAK_SERVICE_CLIENT_SECRET", "STUDENT_SECRET")

        # Fallback: test user credentials for local dev
        self.test_username = getattr(self.settings, "KEYCLOAK_TEST_USERNAME", "student")
        self.test_password = getattr(self.settings, "KEYCLOAK_TEST_PASSWORD", "student")

    @property
    def token_url(self) -> str:
        """Keycloak token endpoint URL."""
        return f"{self.server_url}/realms/{self.realm}/protocol/openid-connect/token"

    async def get_access_token(self) -> str:
        """Get a valid access token, using cache if available.

        Tries in order:
        1. Return cached token if still valid
        2. Client credentials grant (service account)
        3. Password grant (test user) - fallback for development

        Returns
        -------
        str
            Valid JWT access token

        Raises
        ------
        httpx.HTTPStatusError
            If token request fails
        """
        # Return cached token if valid
        if self._cache.is_valid():
            logger.debug("Using cached access token")
            return self._cache.access_token

        # Try client credentials first
        try:
            token_data = await self._request_token_client_credentials()
            self._update_cache(token_data)
            logger.info("Obtained access token via client credentials grant")
            return self._cache.access_token
        except httpx.HTTPStatusError as e:
            if e.response.status_code == 401:
                logger.warning("Client credentials failed, trying password grant: %s", e)
            else:
                raise

        # Fallback to password grant (for development)
        token_data = await self._request_token_password()
        self._update_cache(token_data)
        logger.info("Obtained access token via password grant")
        return self._cache.access_token

    async def _request_token_client_credentials(self) -> dict[str, Any]:
        """Request token using client credentials grant."""
        async with httpx.AsyncClient(timeout=10) as client:
            response = await client.post(
                self.token_url,
                data={
                    "grant_type": "client_credentials",
                    "client_id": self.client_id,
                    "client_secret": self.client_secret,
                },
            )
            response.raise_for_status()
            return response.json()

    async def _request_token_password(self) -> dict[str, Any]:
        """Request token using resource owner password grant."""
        async with httpx.AsyncClient(timeout=10) as client:
            response = await client.post(
                self.token_url,
                data={
                    "grant_type": "password",
                    "client_id": self.client_id,
                    "client_secret": self.client_secret,
                    "username": self.test_username,
                    "password": self.test_password,
                },
            )
            response.raise_for_status()
            return response.json()

    def _update_cache(self, token_data: dict[str, Any]) -> None:
        """Update token cache with new token data."""
        self._cache.access_token = token_data["access_token"]
        # expires_in is in seconds
        expires_in = token_data.get("expires_in", 300)
        self._cache.expires_at = time.time() + expires_in

    async def get_auth_headers(self) -> dict[str, str]:
        """Get authorization headers with valid token.

        Convenience method that returns headers dict ready for httpx.

        Returns
        -------
        dict[str, str]
            Headers dict with Authorization: Bearer <token>
        """
        token = await self.get_access_token()
        return {"Authorization": f"Bearer {token}"}


# Singleton instance for convenience
_auth_instance: KeycloakServiceAuth | None = None


def get_keycloak_auth() -> KeycloakServiceAuth:
    """Get singleton KeycloakServiceAuth instance."""
    global _auth_instance
    if _auth_instance is None:
        _auth_instance = KeycloakServiceAuth()
    return _auth_instance


__all__ = ["KeycloakServiceAuth", "get_keycloak_auth"]
