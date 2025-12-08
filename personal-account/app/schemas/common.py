"""Common schemas used across the application."""
from typing import Generic, TypeVar
from pydantic import BaseModel

T = TypeVar("T")


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
