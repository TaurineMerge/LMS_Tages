"""Custom exceptions for the application."""


class AppException(Exception):
    """Base application exception."""
    
    def __init__(self, message: str, status_code: int = 400):
        self.message = message
        self.status_code = status_code
        super().__init__(self.message)


class NotFoundError(AppException):
    """Resource not found exception."""
    
    def __init__(self, resource: str, identifier: str | None = None):
        message = f"{resource} not found"
        if identifier:
            message = f"{resource} with id '{identifier}' not found"
        super().__init__(message, status_code=404)


class ConflictError(AppException):
    """Conflict exception (e.g., duplicate resource)."""
    
    def __init__(self, message: str):
        super().__init__(message, status_code=409)


class ValidationError(AppException):
    """Validation error exception."""
    
    def __init__(self, message: str):
        super().__init__(message, status_code=422)
