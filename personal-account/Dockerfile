# Используем официальный Python образ
FROM python:3.11-slim 

ENV POETRY_HOME="/opt/poetry" \
    POETRY_VIRTUALENVS_IN_PROJECT=true \
    POETRY_NO_INTERACTION=1
# ENV  PYTHONFAULTHANDLER=1 \
#     PYTHONUNBUFFERED=1 \
#     PYTHONHASHSEED=random \
#     PIP_NO_CACHE_DIR=off \
#     PIP_DISABLE_PIP_VERSION_CHECK=on \
#     PIP_DEFAULT_TIMEOUT=100 \
#     # Poetry's configuration:
#     POETRY_NO_INTERACTION=1 \
#     POETRY_VIRTUALENVS_CREATE=false \
#     POETRY_CACHE_DIR='/var/cache/pypoetry' \
#     POETRY_HOME='/.local/bin/poetry' \
#     POETRY_VERSION=2.2.1 

ENV PATH="$POETRY_HOME/bin:$VENV_PATH/bin:$PATH"

RUN apt-get update && apt-get install --no-install-recommends -y curl \
    && curl -sSL https://install.python-poetry.org | python

# Проверяем установку Poetry и обновляем
RUN ${POETRY_HOME}/bin/poetry --version && ${POETRY_HOME}/bin/poetry self update
RUN poetry config --list
# Устанавливаем рабочую директорию
WORKDIR /app

ADD . /app

# Копируем файл зависимостей
COPY pyproject.toml poetry.lock* /app/

# Копируем остальные файлы
COPY . .

# Устанавливаем зависимости
# RUN ${POETRY_HOME}/bin/poetry install --only=main --no-interaction --no-ansi
RUN /opt/poetry/bin/poetry install --only main --no-interaction --no-ansi

# Копируем остальные файлы *после* установки зависимостей для лучшего использования кэша
COPY . .

# Создаем пользователя для безопасности
RUN useradd -m -u 1000 fastapiuser && chown -R fastapiuser:fastapiuser /app
USER fastapiuser

# Открываем порт
EXPOSE 8000

# Запускаем серверß
CMD ["/root/.poetry/bin/poetry", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]