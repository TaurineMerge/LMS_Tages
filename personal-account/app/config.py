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
        self.API_PREFIX: str = os.getenv("API_PREFIX", "/api/v1")

        # ПУБЛИЧНЯ ЧАСТЬ БУДЕТ ПОТОМ :)
        self.API_PUBLIC_URL: str = os.getenv("API_PUBLIC_URL", "http://localhost:8004")

        # Observability / OpenTelemetry - Core
        self.OTEL_EXPORTER_OTLP_ENDPOINT: str = os.getenv(
            "OTEL_EXPORTER_OTLP_ENDPOINT",
            "http://otel-collector:4317",
        )
        self.OTEL_SERVICE_NAME: str = os.getenv(
            "OTEL_SERVICE_NAME",
            "personal-account-api",
        )
        self.OTEL_EXPORTER_OTLP_INSECURE: bool = os.getenv("OTEL_EXPORTER_OTLP_INSECURE", "true").lower() == "true"

        # OpenTelemetry - Sampling
        self.OTEL_SAMPLING_RATE: str = os.getenv("OTEL_TRACES_SAMPLER_ARG", "1.0")  # 1.0 = 100%

        # OpenTelemetry - Span Limits
        self.OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT: str = os.getenv("OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT", "128")
        self.OTEL_SPAN_EVENT_COUNT_LIMIT: str = os.getenv("OTEL_SPAN_EVENT_COUNT_LIMIT", "128")
        self.OTEL_SPAN_LINK_COUNT_LIMIT: str = os.getenv("OTEL_SPAN_LINK_COUNT_LIMIT", "128")
        self.OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT: str = os.getenv("OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT", "4096")

        # OpenTelemetry - Export Configuration
        self.OTEL_EXPORTER_OTLP_TIMEOUT: str = os.getenv("OTEL_EXPORTER_OTLP_TIMEOUT", "30000")  # ms
        self.OTEL_BSP_MAX_EXPORT_BATCH_SIZE: str = os.getenv("OTEL_BSP_MAX_EXPORT_BATCH_SIZE", "512")
        self.OTEL_BSP_SCHEDULE_DELAY: str = os.getenv("OTEL_BSP_SCHEDULE_DELAY", "5000")  # ms
        self.OTEL_BSP_MAX_QUEUE_SIZE: str = os.getenv("OTEL_BSP_MAX_QUEUE_SIZE", "2048")

        # OpenTelemetry - Feature Flags
        self.OTEL_ENABLE_CONSOLE_EXPORTER: bool = os.getenv("OTEL_ENABLE_CONSOLE_EXPORTER", "false").lower() == "true"
        self.OTEL_ENABLE_LOGGING: bool = os.getenv("OTEL_ENABLE_LOGGING", "true").lower() == "true"
        self.OTEL_ENABLE_SQLALCHEMY: bool = os.getenv("OTEL_ENABLE_SQLALCHEMY", "true").lower() == "true"
        self.OTEL_ENABLE_HTTPX: bool = os.getenv("OTEL_ENABLE_HTTPX", "true").lower() == "true"
        self.OTEL_EXCLUDED_URLS: str = os.getenv("OTEL_EXCLUDED_URLS", "/health,/metrics")

        # Service metadata
        self.SERVICE_VERSION: str = os.getenv("SERVICE_VERSION", "1.0.0")
        self.SERVICE_NAMESPACE: str = os.getenv("SERVICE_NAMESPACE", "lms-tages")
        self.SERVICE_INSTANCE_ID: str = os.getenv("SERVICE_INSTANCE_ID", "unknown")
        self.ENVIRONMENT: str = os.getenv("ENVIRONMENT", "production")

        # Keycloak Authentication
        self.KEYCLOAK_SERVER_URL: str = os.getenv("KEYCLOAK_SERVER_URL", "http://keycloak:8080/auth/")
        self.KEYCLOAK_PUBLIC_URL: str = os.getenv("KEYCLOAK_PUBLIC_URL", "http://localhost:80/auth")
        self.KEYCLOAK_REALM: str = os.getenv("KEYCLOAK_REALM", "student")
        self.KEYCLOAK_ADMIN_REALM: str = os.getenv("KEYCLOAK_ADMIN_REALM", "teacher")
        self.KEYCLOAK_CLIENT_ID: str = os.getenv("KEYCLOAK_CLIENT_ID", "personal-account-client")
        self.KEYCLOAK_CLIENT_SECRET: str = os.getenv("KEYCLOAK_CLIENT_SECRET", "secret")
        self.KEYCLOAK_REDIRECT_URI: str = os.getenv("KEYCLOAK_REDIRECT_URI", "http://localhost/account/callback")

        # Keycloak Admin credentials (for user registration)
        self.KEYCLOAK_ADMIN_USERNAME: str = os.getenv("KEYCLOAK_ADMIN_USERNAME", "admin")
        self.KEYCLOAK_ADMIN_PASSWORD: str = os.getenv("KEYCLOAK_ADMIN_PASSWORD", "admin")

        self.KEYCLOAK_DEFAULT_SCOPE: str = os.getenv("KEYCLOAK_DEFAULT_SCOPE", "openid profile email")
        self.KEYCLOAK_USER_EMAIL_VERIFIED_DEFAULT: bool = (
            os.getenv("KEYCLOAK_USER_EMAIL_VERIFIED_DEFAULT", "true").lower() == "true"
        )

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
