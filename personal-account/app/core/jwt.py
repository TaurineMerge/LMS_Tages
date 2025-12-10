"""JWT helper: decode/validate (JWKS) and optional encode for internal tokens."""
from typing import Optional, Dict, Any, List
import time
import logging

import httpx
from jose import jwt
from jose.exceptions import JWTError, ExpiredSignatureError

logger = logging.getLogger(__name__)


class JwtService:
    def __init__(self, keycloak_url: str, realm: str):
        self.keycloak_url = keycloak_url.rstrip("/")
        self.realm = realm
        self._jwks_cache: Optional[Dict[str, Any]] = None
        self._jwks_time: Optional[float] = None
        self._jwks_ttl = 3600  # seconds

    async def _fetch_jwks(self) -> Dict[str, Any]:
        if self._jwks_cache and self._jwks_time and (time.time() - self._jwks_time) < self._jwks_ttl:
            return self._jwks_cache

        url = f"{self.keycloak_url}/realms/{self.realm}/protocol/openid-connect/certs"
        async with httpx.AsyncClient(timeout=10) as client:
            r = await client.get(url)
            r.raise_for_status()
            self._jwks_cache = r.json()
            self._jwks_time = time.time()
            logger.debug("Fetched JWKS (%d keys)", len(self._jwks_cache.get("keys", [])))
            return self._jwks_cache

    async def decode(self, token: str, audience: Optional[str] = None, issuer: Optional[str] = None) -> Dict[str, Any]:
        """
        Validate and decode token using JWKS. Returns payload dict or raises JWTError/ExpiredSignatureError.
        """
        jwks = await self._fetch_jwks()
        # jose accepts jwks directly as key (works like before)
        return jwt.decode(
            token,
            jwks,
            algorithms=["RS256"],
            audience=audience,
            issuer=issuer,
            options={"verify_exp": True}
        )

    def encode(self, payload: Dict[str, Any], private_key_pem: str, algorithm: str = "RS256", expires_in: Optional[int] = None) -> str:
        """
        Encode internal token (only if you control signing key). Use this only for internal tokens,
        not to fake Keycloak tokens. private_key_pem must be stored securely.
        """
        data = payload.copy()
        now = int(time.time())
        if expires_in:
            data.setdefault("iat", now)
            data.setdefault("exp", now + expires_in)
        return jwt.encode(data, private_key_pem, algorithm=algorithm)