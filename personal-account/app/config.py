"""Application configuration."""
import os
from functools import lru_cache


class settings:
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

        # Observability / OpenTelemetry
        self.OTEL_EXPORTER_OTLP_ENDPOINT: str = os.getenv(
            "OTEL_EXPORTER_OTLP_ENDPOINT",
            "http://otel-collector:4317",
        )
        self.OTEL_SERVICE_NAME: str = os.getenv(
            "OTEL_SERVICE_NAME",
            "personal-account-api",
        )
        self.OTEL_EXPORTER_OTLP_INSECURE: bool = (
            os.getenv("OTEL_EXPORTER_OTLP_INSECURE", "true").lower() == "true"
        )

        # Keycloak Authentication
        self.KEYCLOAK_SERVER_URL: str = os.getenv("KEYCLOAK_SERVER_URL", "http://keycloak:8080")
        self.KEYCLOAK_PUBLIC_URL: str = os.getenv("KEYCLOAK_PUBLIC_URL", "http://localhost:8080")
        self.KEYCLOAK_REALM: str = os.getenv("KEYCLOAK_REALM", "student")
        self.KEYCLOAK_CLIENT_ID: str = os.getenv("KEYCLOAK_CLIENT_ID", "personal-account-client")
        self.KEYCLOAK_CLIENT_SECRET: str = os.getenv("KEYCLOAK_CLIENT_SECRET", "secret")
        self.KEYCLOAK_REDIRECT_URI: str = os.getenv("KEYCLOAK_REDIRECT_URI", "http://localhost/account/callback")
        
        # Keycloak Admin credentials (for user registration)
        self.KEYCLOAK_ADMIN_USERNAME: str = os.getenv("KEYCLOAK_ADMIN_USERNAME", "admin")
        self.KEYCLOAK_ADMIN_PASSWORD: str = os.getenv("KEYCLOAK_ADMIN_PASSWORD", "admin")


    @property
    def database_url(self) -> str:
        """Construct database URL."""
        return (
            f"postgresql://{self.DATABASE_USER}:{self.DATABASE_PASSWORD}"
            f"@{self.DATABASE_HOST}:{self.DATABASE_PORT}/{self.DATABASE_NAME}"
        )


@lru_cache()
def get_settings() -> settings:
    """Get cached settings instance."""
    return settings()
