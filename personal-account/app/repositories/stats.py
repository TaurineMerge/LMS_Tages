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
from app.redis_client import get_redis_client
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
    async def get_unproceMARK_RAW_ATTEMPT_PROCESSEDssed_user_stats(self, limit: int = 100) -> list[dict[str, Any]]:
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

    @traced("stats.get_cached_user_stats", record_args=True, record_result=True)
    async def get_cached_user_stats(self, student_id: UUID) -> dict[str, Any] | None:
        """Get cached user statistics from Redis.

        Returns aggregated statistics for the user if available in cache,
        otherwise returns None.
        """
        redis_client = get_redis_client()
        cache_key = f"user_stats:{student_id}"
        cached_data = await redis_client.get(cache_key)
        if cached_data:
            return json.loads(cached_data)
        return None

    @traced("stats.set_cached_user_stats", record_args=True, record_result=True)
    async def set_cached_user_stats(self, student_id: UUID, stats: dict[str, Any]) -> None:
        """Cache user statistics in Redis with TTL.

        Stores aggregated statistics for the user with configured TTL.
        """
        from app.config import get_settings

        settings = get_settings()
        redis_client = get_redis_client()
        cache_key = f"user_stats:{student_id}"
        await redis_client.setex(cache_key, settings.REDIS_CACHE_TTL, json.dumps(stats))

    @traced("stats.get_user_stats_from_db", record_args=True, record_result=True)
    async def get_user_stats_from_db(self, student_id: UUID) -> dict[str, Any]:
        """Get aggregated user statistics from database.

        Calculates statistics from the business table for the given user.
        This is a placeholder - implement actual aggregation logic based on requirements.
        """
        # Placeholder: implement actual aggregation query
        # For example: count attempts, average scores, etc.
        result = await fetch_all(
            "SELECT COUNT(*) as total_attempts FROM tests.test_attempt_b WHERE student_id = $1", [str(student_id)]
        )
        return {
            "student_id": str(student_id),
            "total_attempts": result[0]["total_attempts"] if result else 0,
            # Add more aggregated fields as needed
        }

    @traced("stats.invalidate_user_stats_cache", record_args=True, record_result=True)
    async def invalidate_user_stats_cache(self, student_id: UUID) -> None:
        """Invalidate cached user statistics in Redis.

        Removes the cached statistics for the user, forcing fresh calculation on next access.
        """
        redis_client = get_redis_client()
        cache_key = f"user_stats:{student_id}"
        await redis_client.delete(cache_key)

    @traced("stats.get_user_stats", record_args=True, record_result=True)
    async def get_user_stats(self, student_id: UUID) -> dict[str, Any]:
        """Get user statistics with caching.

        First checks Redis cache, if not found calculates from database and caches the result.
        """
        # Try cache first
        cached_stats = await self.get_cached_user_stats(student_id)
        if cached_stats:
            return cached_stats

        # Calculate from database
        stats = await self.get_user_stats_from_db(student_id)

        # Cache the result
        await self.set_cached_user_stats(student_id, stats)

        return stats
