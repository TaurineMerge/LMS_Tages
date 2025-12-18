"""Business repository for stats-related DB operations.

This module contains the database operations that populate the
business table `tests.test_attempt_b` from raw payloads.

Keeping SQL here keeps the processor focused on orchestration and
validation and makes the DB interactions easier to test.
"""

import json
from typing import Any
from uuid import UUID

from app.database import execute, execute_returning, fetch_all
from app.db import queries as q
from app.telemetry import traced


class StatsRepository:
    """Repository that encapsulates business table operations for stats.

    Methods are small wrappers around SQL statements. All public methods
    are instrumented with the `@traced()` decorator so calls show in traces.
    """

    def __init__(self):
        # repository is stateless; DB helpers are used directly
        pass

    @traced("stats.get_unprocessed_attempts", record_args=True, record_result=True)
    async def get_unprocessed_attempts(self, limit: int = 100) -> list[dict[str, Any]]:
        return await fetch_all(q.SELECT_UNPROCESSED_ATTEMPTS, {"limit": limit})

    @traced("stats.mark_attempt_processed", record_args=True, record_result=True)
    async def mark_attempt_processed(self, raw_id: UUID, error: str | None) -> None:
        # Record processing completion. The centralized query marks processed
        # and updates the timestamp. If storing error strings is required,
        # add a dedicated query to app/db/queries.py and use it here.
        await execute_returning(q.MARK_RAW_ATTEMPT_PROCESSED, {"id": str(raw_id)})

    @traced("stats.get_unprocessed_user_stats", record_args=True, record_result=True)
    async def get_unprocessed_user_stats(self, limit: int = 100) -> list[dict[str, Any]]:
        return await fetch_all(q.SELECT_UNPROCESSED_USER_STATS, {"limit": limit})

    @traced("stats.mark_user_stats_processed", record_args=True, record_result=True)
    async def mark_user_stats_processed(self, raw_id: UUID, error: str | None) -> None:
        await execute_returning(q.MARK_RAW_USER_STATS_PROCESSED, {"id": str(raw_id)})

    @traced("stats.upsert_attempt", record_args=True, record_result=True)
    async def upsert_attempt(self, payload: dict[str, Any]) -> None:
        """Insert or update a record in tests.test_attempt_b from validated payload.

        Expects a payload already validated by ContractManager. Accepts the
        raw JSON structure produced by the testing service and maps fields
        into the business table.
        """
        params = {
            "id": str(UUID(payload["attempt_id"])),
            "student_id": str(UUID(payload["student_id"])),
            "test_id": str(UUID(payload["test_id"])) if payload.get("test_id") else None,
            "date_of_attempt": payload.get("date_of_attempt"),
            "point": payload.get("point"),
            # Serialize JSON blobs to strings for safe binding
            "result": json.dumps(payload.get("result", {})),
            "completed": payload.get("completed", False),
            "passed": payload.get("passed"),
            "certificate_id": str(UUID(payload["certificate_id"])) if payload.get("certificate_id") else None,
            "snapshot": payload.get("attempt_snapshot_s3"),
            "version": json.dumps(payload.get("attempt_version", {})) if payload.get("attempt_version") else None,
            "meta": json.dumps(payload.get("meta", {})),
        }

        await execute(q.TEST_ATTEMPT_UPSERT, params)
