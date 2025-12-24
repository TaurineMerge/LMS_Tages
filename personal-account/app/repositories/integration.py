"""Integration repository — stores raw data received from external services.

This repository writes raw payloads into the `integration` schema tables
created by `init-sql/migrate-002.sql`:

- `integration.raw_user_stats`
- `integration.raw_attempts`

It uses the same style as other repositories in `app/repositories` and is
instrumented with the `@traced()` decorator from `app.telemetry`.

Public methods
--------------
- `insert_raw_user_stats(student_id, payload)` — insert a raw user stats JSON.
- `insert_raw_attempt(external_attempt_id, student_id, test_id, payload)` — insert/update raw attempt.
- `get_unprocessed_user_stats(limit)` — fetch unprocessed user stats records.
- `get_unprocessed_attempts(limit)` — fetch unprocessed attempts.
- `mark_user_stats_processed(id)` / `mark_attempt_processed(id)` — mark processed.
- `increment_attempt_processing(id)` — increment processing attempts counter.

Example
-------
```py
from app.repositories.integration import integration_repository

await integration_repository.insert_raw_user_stats(student_id, payload)
```"""

from __future__ import annotations

import json
from typing import Any
from uuid import UUID

from app.database import execute_returning, fetch_all, fetch_one

# Delay importing telemetry and database helpers so docs generation or
# lightweight introspection doesn't require opentelemetry / DB deps to be
# installed. We prefer to import these at runtime inside methods.
from app.db import queries as q
from app.repositories.base import base_repository

try:  # pragma: no cover - best-effort import for runtime
    from app.telemetry import traced
except ImportError:  # fallback to no-op decorator when telemetry is missing

    def traced(*_args, **_kwargs):
        def _decorator(func):
            return func

        return _decorator


class integration_repository(base_repository):
    """Repository to manage raw integration data.

    The repository writes directly into the `integration` schema. It follows
    the same interface and naming conventions as other repositories in the
    project.
    """

    def __init__(self):
        super().__init__(
            table_name="raw_user_stats",
            schema="integration",
            default_order_by="received_at",
            orderable_columns={"received_at", "created_at", "updated_at"},
        )

    @traced()
    async def insert_raw_user_stats(self, student_id: UUID, payload: dict[str, Any]) -> dict[str, Any] | None:
        """Insert a raw user stats payload.

        Parameters
        ----------
        student_id: UUID
            Local student id to associate the payload with.
        payload: dict
            Raw JSON payload received from testing service.
        """

        params = {
            "student_id": student_id,
            "payload": json.dumps(payload),
            "received_at": None,
            "processed": False,
            "error_message": None,
        }
        return await execute_returning(q.INSERT_RAW_USER_STATS, params)

    @traced()
    async def insert_raw_attempt(
        self, external_attempt_id: UUID, student_id: UUID, test_id: UUID | None, payload: dict[str, Any]
    ) -> dict[str, Any] | None:
        """Insert or update a raw attempt record.

        On conflict by `external_attempt_id` the payload and received_at are
        updated (we rely on SQL ON CONFLICT clause in the query).
        """
        # Local import to avoid importing DB layer at module import time

        params = {
            "external_attempt_id": external_attempt_id,
            "student_id": student_id,
            "test_id": test_id,
            "payload": json.dumps(payload),
            "received_at": None,
            "processed": False,
            "processing_attempts": 0,
            "error_message": None,
        }
        return await execute_returning(q.INSERT_RAW_ATTEMPT, params)

    @traced()
    async def get_unprocessed_user_stats(self, limit: int = 100) -> list[dict[str, Any]]:
        """Return unprocessed raw user stats records ordered by received_at."""
        return await fetch_all(q.SELECT_UNPROCESSED_USER_STATS, {"limit": limit})

    @traced()
    async def get_unprocessed_attempts(self, limit: int = 100) -> list[dict[str, Any]]:
        """Return unprocessed raw attempts records ordered by received_at."""

        return await fetch_all(q.SELECT_UNPROCESSED_ATTEMPTS, {"limit": limit})

    @traced()
    async def mark_user_stats_processed(self, record_id: UUID) -> dict[str, Any] | None:
        """Mark raw_user_stats record as processed."""

        return await execute_returning(q.MARK_RAW_USER_STATS_PROCESSED, {"id": record_id})

    @traced()
    async def mark_attempt_processed(self, record_id: UUID) -> dict[str, Any] | None:
        """Mark raw_attempts record as processed."""

        return await execute_returning(q.MARK_RAW_ATTEMPT_PROCESSED, {"id": record_id})

    @traced()
    async def increment_attempt_processing(self, record_id: UUID) -> int | None:
        """Increment processing_attempts counter and return new value."""

        result = await fetch_one(q.INCREMENT_RAW_ATTEMPT_PROCESSING, {"id": record_id})
        return result["processing_attempts"] if result else None


# Singleton instance
integration_repository = integration_repository()
