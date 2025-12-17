FROM python:3.11-slim AS builder

WORKDIR /app

# Установка системных зависимостей, необходимых для некоторых пакетов Python
RUN apt-get update \
    && apt-get install --no-install-recommends -y \
        build-essential \
        gcc \
    && rm -rf /var/lib/apt/lists/*

# Установка Poetry напрямую через pip (более надёжный способ для Docker)
RUN pip install poetry

# Установка переменных окружения Poetry до выполнения команд
ENV POETRY_NO_INTERACTION=1 \
    POETRY_VIRTUALS_CREATE=false \
    POETRY_CACHE_DIR=/tmp/poetry_cache

# Проверка установки Poetry и местоположения
RUN which poetry || echo "Poetry executable not found in PATH"
RUN poetry --version

# Copy only dependency files first (for layer caching)
COPY pyproject.toml poetry.lock* /app/

# Install dependencies in system Python (skip virtualenv for Docker)
RUN poetry config virtualenvs.create false \
    && poetry install --no-root --only main --no-interaction --no-ansi

# Clean up Poetry cache to reduce layer size
RUN rm -rf $POETRY_CACHE_DIR

# Copy application code (won't invalidate dependency layer)
COPY . /app

# Stage 2: Runtime image (minimal)
FROM python:3.11-slim

WORKDIR /app

# Copy installed packages and app code
COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /app /app

RUN useradd -m -u 1000 fastapiuser && chown -R fastapiuser:fastapiuser /app
USER fastapiuser

EXPOSE 8000

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
