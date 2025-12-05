"""Application configuration."""
import os
from functools import lru_cache


class Settings:
    """Application settings loaded from environment variables."""
    
    def __init__(self):
        self.DATABASE_HOST: str = os.getenv("DATABASE_HOST", "app-db")
        self.DATABASE_PORT: int = int(os.getenv("DATABASE_PORT", "5432"))
        self.DATABASE_NAME: str = os.getenv("DATABASE_NAME", "appdb")
        self.DATABASE_USER: str = os.getenv("DATABASE_USER", "appuser")
        self.DATABASE_PASSWORD: str = os.getenv("DATABASE_PASSWORD", "password")
        self.DATABASE_POOL_MIN_SIZE: int = int(os.getenv("DATABASE_POOL_MIN_SIZE", "5"))
        self.DATABASE_POOL_MAX_SIZE: int = int(os.getenv("DATABASE_POOL_MAX_SIZE", "20"))
        
        self.DEBUG: bool = os.getenv("DEBUG", "false").lower() == "true"
        self.API_PREFIX: str = "/api/v1"

    @property
    def database_url(self) -> str:
        """Construct database URL."""
        return (
            f"postgresql://{self.DATABASE_USER}:{self.DATABASE_PASSWORD}"
            f"@{self.DATABASE_HOST}:{self.DATABASE_PORT}/{self.DATABASE_NAME}"
        )


@lru_cache()
def get_settings() -> Settings:
    """Get cached settings instance."""
    return Settings()
