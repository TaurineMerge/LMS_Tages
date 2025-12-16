# Stage 1: Build dependencies layer (cached unless pyproject.toml/poetry.lock change)
FROM python:3.11-slim AS builder

WORKDIR /app

# Install Poetry once and cache this layer
RUN apt-get update \
    && apt-get install --no-install-recommends -y curl \
    && curl -sSL https://install.python-poetry.org | python3 - \
    && rm -rf /var/lib/apt/lists/*

    ENV PATH="/root/.local/bin:$PATH"

# Copy only dependency files first (for layer caching)
COPY pyproject.toml poetry.lock* /app/

# Install dependencies in system Python (skip virtualenv for Docker)
RUN poetry config virtualenvs.create false \
    && poetry install --no-root --only main --no-interaction --no-ansi

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