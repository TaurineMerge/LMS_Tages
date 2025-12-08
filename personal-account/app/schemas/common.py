"""Common schemas used across the application."""
from typing import Any, Generic, TypeVar
from pydantic import BaseModel

T = TypeVar("T")


class api_response(BaseModel, Generic[T]):
    """Standard API response wrapper."""
    
    status: str = "success"
    data: T | None = None
    message: str | None = None


class api_error_response(BaseModel):
    """Standard API error response."""
    
    status: str = "error"
    error: str
    details: dict[str, Any] | None = None


class paginated_response(BaseModel, Generic[T]):
    """Paginated response schema."""
    
    data: list[T]
    total: int
    page: int
    limit: int


class error_response(BaseModel):
    """Error response schema."""
    
    error: str


class message_response(BaseModel):
    """Simple message response."""
    
    message: str


class token_response(BaseModel):
    """Token response from Keycloak."""
    
    access_token: str
    refresh_token: str | None = None
    token_type: str = "Bearer"
    expires_in: int
    refresh_expires_in: int | None = None
    scope: str | None = None


class user_info_response(BaseModel):
    """User info from token."""
    
    sub: str
    email: str | None = None
    preferred_username: str | None = None
    name: str | None = None
    given_name: str | None = None
    family_name: str | None = None
    email_verified: bool | None = None


class register_request(BaseModel):
    """User registration request."""
    
    username: str
    email: str
    password: str
    first_name: str | None = None
    last_name: str | None = None


class register_response(BaseModel):
    """User registration response."""
    
    user_id: str
    username: str
    email: str
    message: str = "Registration successful"

