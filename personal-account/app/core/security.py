"""JWT security and token validation."""

import logging
from typing import Optional
from datetime import datetime
import httpx

from fastapi import HTTPException, Depends, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jose import jwt
from jose.exceptions import JWTError, ExpiredSignatureError
from pydantic import BaseModel, Field

from app.config import get_settings
from app.telemetry import traced

logger = logging.getLogger(__name__)
settings = get_settings()

# Security scheme for FastAPI
security = HTTPBearer()


class TokenPayload(BaseModel):
    """JWT token payload from Keycloak."""
    sub: str = Field(..., description="User ID (subject)")
    email: Optional[str] = None
    name: Optional[str] = None
    preferred_username: Optional[str] = None
    exp: int = Field(..., description="Token expiration time")
    iat: Optional[int] = None
    
    class Config:
        json_schema_extra = {
            "example": {
                "sub": "user-id-uuid",
                "email": "user@example.com",
                "name": "User Name",
                "preferred_username": "username",
                "exp": 1765284914,
                "iat": 1765284614
            }
        }


class JWTValidator:
    """Validate JWT tokens using Keycloak JWKS."""
    
    def __init__(
        self,
        keycloak_url: str,
        realm: str,
        client_id: str,
    ):
        self.keycloak_url = keycloak_url.rstrip("/")
        self.realm = realm
        self.client_id = client_id
        self.issuer = f"{self.keycloak_url}/realms/{self.realm}"
        self._jwks_cache = None
        self._cache_time = None
    
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
            # Получаем JWKS для валидации подписи
            jwks = await self._get_jwks()
            
            # Валидируем токен через python-jose
            payload = jwt.decode(
                token,
                jwks,
                algorithms=["RS256"],
                audience="account",
                issuer=self.issuer,
                options={"verify_exp": True}
            )
            
            # Преобразуем в типизированный объект
            token_payload = TokenPayload(**payload)
            
            logger.info(
                f"Token validated successfully for user {token_payload.sub}"
            )
            return token_payload
            
        except ExpiredSignatureError:
            logger.warning("JWT token has expired")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Token has expired",
                headers={"WWW-Authenticate": "Bearer"}
            )
        except JWTError as e:
            logger.warning(f"JWT validation failed: {e}")
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid token",
                headers={"WWW-Authenticate": "Bearer"}
            )
        except Exception as e:
            logger.error(f"Unexpected error during token validation: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Token validation error"
            )
    
    @traced("jwt_validator.get_jwks", record_args=True, record_result=True)
    async def _get_jwks(self) -> dict:
        """
        Fetch and cache JWKS from Keycloak.
        
        JWKS (JSON Web Key Set) содержит публичные ключи
        для верификации JWT подписей.
        """
        import time
        
        # Проверяем кэш (1 час = 3600 секунд)
        if self._jwks_cache and self._cache_time:
            if time.time() - self._cache_time < 3600:
                logger.debug("Using cached JWKS")
                return self._jwks_cache
        
        try:
            jwks_url = (
                f"{self.keycloak_url}/realms/{self.realm}"
                "/protocol/openid-connect/certs"
            )
            
            logger.debug(f"Fetching JWKS from {jwks_url}")
            
            async with httpx.AsyncClient(timeout=10) as client:
                response = await client.get(jwks_url)
                response.raise_for_status()
                
                self._jwks_cache = response.json()
                self._cache_time = time.time()
                
                logger.debug(f"JWKS fetched successfully, {len(self._jwks_cache.get('keys', []))} keys available")
                
                return self._jwks_cache
                
        except httpx.HTTPError as e:
            logger.error(f"Failed to fetch JWKS from {jwks_url}: {e}")
            raise HTTPException(
                status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
                detail="Could not verify token"
            )


# Создаём singleton validator
_jwt_validator = JWTValidator(
    keycloak_url=settings.KEYCLOAK_SERVER_URL,
    realm=settings.KEYCLOAK_REALM,
    client_id=settings.KEYCLOAK_CLIENT_ID,
)


@traced("security.get_current_user", record_result=False)
async def get_current_user(
    credentials: HTTPAuthorizationCredentials = Depends(security),
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
    if not credentials:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Not authenticated",
            headers={"WWW-Authenticate": "Bearer"}
        )
    
    return await _jwt_validator.validate_token(credentials.credentials)


@traced("security.get_current_user_optional", record_result=False)
async def get_current_user_optional(
    credentials: Optional[HTTPAuthorizationCredentials] = Depends(security),
) -> Optional[TokenPayload]:
    """
    Get current user if authenticated, otherwise return None.
    
    Используется для опциональной аутентификации.
    """
    if not credentials:
        return None
    
    try:
        return await _jwt_validator.validate_token(credentials.credentials)
    except HTTPException:
        # Если токен невалидный, просто возвращаем None
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
    async def check_roles(
        current_user: TokenPayload = Depends(get_current_user)
    ) -> None:
        # В этом примере просто проверяем что пользователь аутентифицирован
        # Можно расширить для проверки realm roles или client roles
        if not current_user:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Insufficient permissions"
            )
    
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
