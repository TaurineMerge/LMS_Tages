"""Common schemas used across the application."""
from typing import Generic, TypeVar
from pydantic import BaseModel

T = TypeVar("T")


class PaginatedResponse(BaseModel, Generic[T]):
    """Paginated response schema."""
    
    data: list[T]
    total: int
    page: int
    limit: int


class ErrorResponse(BaseModel):
    """Error response schema."""
    
    error: str


class MessageResponse(BaseModel):
    """Simple message response."""
    
    message: str
