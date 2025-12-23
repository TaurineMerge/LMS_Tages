"""Application configuration."""

import logging
from functools import lru_cache

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    model_config = SettingsConfigDict(env_file=(".env", ".env.local"), env_file_encoding="utf-8", extra="ignore")

    DATABASE_HOST: str = "app-db"
    DATABASE_PORT: int = 5432
    DATABASE_NAME: str = "appdb"
    DATABASE_USER: str = "appuser"
    DATABASE_PASSWORD: str = "password"
    DATABASE_POOL_MIN_SIZE: int = 5
    DATABASE_POOL_MAX_SIZE: int = 20

    DEBUG: bool = False

    API_PREFIX: str = "/api/v1"

    # ПУБЛИЧНЯ ЧАСТЬ БУДЕТ ПОТОМ :)
    API_PUBLIC_URL: str = "http://localhost:8004"

    # Observability / OpenTelemetry - Core
    OTEL_EXPORTER_OTLP_ENDPOINT: str = "http://otel-collector:4317"
    OTEL_SERVICE_NAME: str = "personal-account-api"
    OTEL_EXPORTER_OTLP_INSECURE: bool = True

    # OpenTelemetry - Sampling
    OTEL_SAMPLING_RATE: str = "1.0"  # 1.0 = 100%

    # OpenTelemetry - Span Limits
    OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT: str = "128"
    OTEL_SPAN_EVENT_COUNT_LIMIT: str = "128"
    OTEL_SPAN_LINK_COUNT_LIMIT: str = "128"
    OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT: str = "4096"

    # OpenTelemetry - Export Configuration
    OTEL_EXPORTER_OTLP_TIMEOUT: str = "30000"  # ms
    OTEL_BSP_MAX_EXPORT_BATCH_SIZE: str = "512"
    OTEL_BSP_SCHEDULE_DELAY: str = "5000"  # ms
    OTEL_BSP_MAX_QUEUE_SIZE: str = "2048"

    # OpenTelemetry - Feature Flags
    OTEL_ENABLE_CONSOLE_EXPORTER: bool = False
    OTEL_ENABLE_LOGGING: bool = True
    OTEL_ENABLE_SQLALCHEMY: bool = True
    OTEL_ENABLE_HTTPX: bool = True
    OTEL_EXCLUDED_URLS: str = "/health,/metrics"

    # Service metadata
    SERVICE_VERSION: str = "1.0.0"
    SERVICE_NAMESPACE: str = "lms-tages"
    SERVICE_INSTANCE_ID: str = "unknown"
    ENVIRONMENT: str = "production"

    # Keycloak Authentication
    KEYCLOAK_SERVER_URL: str = "http://keycloak:8080"
    KEYCLOAK_PUBLIC_URL: str = "http://localhost:8080"
    KEYCLOAK_REALM: str = "student"
    KEYCLOAK_ADMIN_REALM: str = "teacher"
    KEYCLOAK_CLIENT_ID: str = "personal-account-client"
    KEYCLOAK_CLIENT_SECRET: str = "secret"
    KEYCLOAK_REDIRECT_URI: str = "http://localhost/account/callback"

    # Keycloak Admin credentials (for user registration)
    KEYCLOAK_ADMIN_USERNAME: str = "admin"
    KEYCLOAK_ADMIN_PASSWORD: str = "admin"

    KEYCLOAK_DEFAULT_SCOPE: str = "openid profile email"
    KEYCLOAK_USER_EMAIL_VERIFIED_DEFAULT: bool = True

    # Redis Configuration
    REDIS_HOST: str = "redis"
    REDIS_PORT: int = 6379
    REDIS_DB: int = 0
    REDIS_PASSWORD: str | None = None
    REDIS_CACHE_TTL: int = 1800  # 30 minutes in seconds

    # Stats Worker intervals (in seconds)
    STATS_WORKER_FETCH_INTERVAL: int = 60
    STATS_WORKER_PROCESS_INTERVAL: int = 15

    # Testing configuration
    TESTING_BASE_URL: str = "http://testing:8085"

    def model_post_init(self, __context):
        """Override defaults based on ENVIRONMENT."""
        if self.ENVIRONMENT == "local":
            logging.info("Applying local environment settings")
            self.STATS_WORKER_FETCH_INTERVAL = 10  # 10 сек для тестирования
            self.STATS_WORKER_PROCESS_INTERVAL = 15  # 15 сек
            # Локальная разработка: Python локально, сервисы на localhost (порты из docker-compose-dev-py.yml)
            self.DATABASE_HOST = "localhost"
            self.REDIS_HOST = "localhost"
            self.KEYCLOAK_SERVER_URL = "http://localhost:8080"  # Убрал /auth - KeycloakOpenID добавит сам
            self.KEYCLOAK_PUBLIC_URL = "http://localhost:8080"  # Оставил /auth для браузера
            self.KEYCLOAK_REDIRECT_URI = "http://localhost:8000/callback"
            self.OTEL_EXPORTER_OTLP_ENDPOINT = "http://localhost:4317"
            self.TESTING_BASE_URL = "http://localhost:8085"
        # Для production/development оставляем defaults (Docker networks)

    @property
    def testing_base_url(self) -> str:
        """Get testing service URL based on environment."""
        if self.ENVIRONMENT == "local":
            return "http://localhost:8085"
        return self.TESTING_BASE_URL

    @property
    def database_url(self) -> str:
        """Construct database URL."""
        base_url = (
            f"postgresql://{self.DATABASE_USER}:{self.DATABASE_PASSWORD}"
            f"@{self.DATABASE_HOST}:{self.DATABASE_PORT}/{self.DATABASE_NAME}"
        )
        # Disable SSL for local development (asyncpg defaults to SSL)
        if self.is_local:
            return f"{base_url}?ssl=disable"
        return base_url

    @property
    def is_local(self) -> bool:
        """Check if running in local development mode."""
        return self.ENVIRONMENT == "local"

    @property
    def root_path(self) -> str:
        """Root path based on environment. Empty for local, /account for docker."""
        return "" if self.is_local else "/account"

    @property
    def url_prefix(self) -> str:
        """URL prefix for templates. Empty for local, /account for docker."""
        return "" if self.is_local else "/account"


@lru_cache
def get_settings() -> Settings:
    """Get cached settings instance."""
    return Settings()
