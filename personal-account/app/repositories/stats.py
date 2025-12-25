"""Business repository for stats-related DB operations.

This module contains the database operations that populate the
business table `tests.test_attempt_b` from raw payloads and manage
aggregated statistics with Redis caching.
"""

import json
import logging
from datetime import date, datetime
from typing import Any
from uuid import UUID

from app.config import get_settings
from app.database import execute, execute_returning, fetch_all, fetch_one
from app.db import queries as q
from app.redis_client import get_redis_client
from app.telemetry import traced

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)


def _convert_uuids_to_strings(obj: Any) -> Any:
    """Recursively convert UUID objects to strings for JSON serialization.

    Args:
        obj: Object to convert

    Returns:
        Object with UUIDs converted to strings
    """
    if isinstance(obj, UUID):
        return str(obj)
    elif isinstance(obj, (date, datetime)):
        return obj.isoformat()
    elif isinstance(obj, dict):
        return {key: _convert_uuids_to_strings(value) for key, value in obj.items()}
    elif isinstance(obj, list):
        return [_convert_uuids_to_strings(item) for item in obj]
    else:
        return obj


class DateTimeEncoder(json.JSONEncoder):
    """Custom JSON encoder that handles date/datetime objects."""

    def default(self, obj):
        if isinstance(obj, (date, datetime)):
            return obj.isoformat()
        return super().default(obj)


class StatsRepository:
    """Repository for stats operations with Redis caching.

    Handles raw attempt processing, aggregated stats calculation,
    and Redis cache management.
    """

    def __init__(self):
        self._settings = get_settings()

    @traced("stats.get_unprocessed_attempts", record_args=True, record_result=True)
    async def get_unprocessed_attempts(self, limit: int = 100) -> list[dict[str, Any]]:
        """Get unprocessed raw attempts from integration table."""
        return await fetch_all(q.SELECT_UNPROCESSED_ATTEMPTS, {"limit": limit})

    @traced("stats.mark_attempt_processed", record_args=True, record_result=True)
    async def mark_attempt_processed(self, raw_id: UUID, error: str | None) -> None:
        """Mark a raw attempt as processed."""
        await execute_returning(q.MARK_RAW_ATTEMPT_PROCESSED, {"id": str(raw_id)})

    @traced("stats.get_unprocessed_user_stats", record_args=True, record_result=True)
    async def get_unprocessed_user_stats(self, limit: int = 100) -> list[dict[str, Any]]:
        """Get unprocessed raw user stats from integration table."""
        return await fetch_all(q.SELECT_UNPROCESSED_USER_STATS, {"limit": limit})

    @traced("stats.mark_user_stats_processed", record_args=True, record_result=True)
    async def mark_user_stats_processed(self, raw_id: UUID, error: str | None) -> None:
        """Mark raw user stats as processed."""
        await execute_returning(q.MARK_RAW_USER_STATS_PROCESSED, {"id": str(raw_id)})

    @traced("stats.upsert_attempt", record_args=True, record_result=True)
    async def upsert_attempt(self, payload: dict[str, Any]) -> None:
        """Upsert a record in tests.test_attempt_b from validated payload."""
        # Convert date_of_attempt to datetime.date if it's a string
        date_of_attempt = payload.get("date_of_attempt")
        if isinstance(date_of_attempt, str):
            date_of_attempt = datetime.fromisoformat(date_of_attempt).date()

        params = {
            "id": str(UUID(payload["attempt_id"])),
            "student_id": str(UUID(payload["student_id"])),
            "test_id": str(UUID(payload["test_id"])) if payload.get("test_id") else None,
            "date_of_attempt": date_of_attempt,  # Now datetime.date
            "point": payload.get("point"),
            "result": json.dumps(payload.get("result", {})),
            "completed": payload.get("completed", False),
            "passed": payload.get("passed"),
            "certificate_id": str(UUID(payload["certificate_id"])) if payload.get("certificate_id") else None,
            "snapshot": payload.get("attempt_snapshot_s3"),
            "version": json.dumps(payload.get("attempt_version", {})) if payload.get("attempt_version") else None,
            "meta": json.dumps(payload.get("meta", {})),
        }
        await execute(q.TEST_ATTEMPT_UPSERT, params)

    @traced("stats.calculate_and_save_aggregated", record_args=True, record_result=True)
    async def calculate_and_save_aggregated(self, student_id: UUID) -> dict[str, Any]:
        """Calculate aggregated stats from DB and save to aggregated table.

        Args:
            student_id: Student UUID

        Returns:
            Aggregated statistics dictionary
        """
        # Calculate from business table
        result = await fetch_one(q.STUDENT_STATS_CALCULATE, {"student_id": str(student_id)})

        if not result:
            stats = {
                "student_id": str(student_id),
                "total_attempts": 0,
                "passed_attempts": 0,
                "failed_attempts": 0,
                "avg_score": 0.0,
                "total_tests_taken": 0,
                "last_attempt_at": None,
            }
        else:
            stats = {
                "student_id": str(student_id),
                "total_attempts": result.get("total_attempts", 0),
                "passed_attempts": result.get("passed_attempts", 0),
                "failed_attempts": result.get("failed_attempts", 0),
                "avg_score": float(result.get("avg_score", 0)),
                "total_tests_taken": result.get("total_tests_taken", 0),
                "last_attempt_at": str(result["last_attempt_at"]) if result.get("last_attempt_at") else None,
            }

        # Convert last_attempt_at to datetime.date for PostgreSQL DATE column
        last_attempt_at_db = (
            datetime.fromisoformat(stats["last_attempt_at"]).date() if stats["last_attempt_at"] else None
        )

        # Save to aggregated table
        params = {
            "student_id": str(student_id),
            "total_attempts": stats["total_attempts"],
            "passed_attempts": stats["passed_attempts"],
            "failed_attempts": stats["failed_attempts"],
            "avg_score": stats["avg_score"],
            "total_tests_taken": stats["total_tests_taken"],
            "last_attempt_at": last_attempt_at_db,  # datetime.date for DB
            "stats_json": json.dumps(stats),  # stats dict with string last_attempt_at
        }
        await execute(q.STUDENT_STATS_AGGREGATED_UPSERT, params)
        logger.info("Saved aggregated stats for student %s", student_id)
        return stats

    @traced("stats.get_from_redis", record_args=True, record_result=True)
    async def get_from_redis(self, student_id: UUID) -> dict[str, Any] | None:
        """Get cached stats from Redis."""
        redis_client = get_redis_client()
        cache_key = f"user_stats:{student_id}"
        cached = await redis_client.get(cache_key)
        if cached:
            stats = json.loads(cached)
            # Ensure last_attempt_at is a string for JSON serialization
            if "last_attempt_at" in stats and isinstance(stats["last_attempt_at"], date):
                stats["last_attempt_at"] = stats["last_attempt_at"].isoformat()
            logger.debug("Cache hit for student %s", student_id)
            return stats
        return None

    @traced("stats.get_cached_full_data", record_args=True, record_result=True)
    async def get_cached_full_data(self, student_id: UUID) -> dict[str, Any] | None:
        """Get cached full data (stats + certificates + attempts + raw data) from Redis."""
        redis_client = get_redis_client()
        cache_key = f"user_full_data:{student_id}"
        cached = await redis_client.get(cache_key)
        if cached:
            data = json.loads(cached)
            logger.debug("Full data cache hit for student %s", student_id)
            return data
        return None

    @traced("stats.save_to_redis", record_args=True, record_result=True)
    async def save_to_redis(self, student_id: UUID, stats: dict[str, Any]) -> None:
        """Save stats to Redis with TTL."""
        redis_client = get_redis_client()
        cache_key = f"user_stats:{student_id}"
        # Convert UUIDs to strings for JSON serialization
        serializable_stats = _convert_uuids_to_strings(stats)
        await redis_client.setex(
            cache_key, self._settings.REDIS_CACHE_TTL, json.dumps(serializable_stats, cls=DateTimeEncoder)
        )
        logger.debug("Cached stats for student %s with TTL %d", student_id, self._settings.REDIS_CACHE_TTL)

    @traced("stats.save_full_data_to_redis", record_args=True, record_result=True)
    async def save_full_data_to_redis(self, student_id: UUID, data: dict[str, Any]) -> None:
        """Save full data (stats + certificates + attempts + raw data) to Redis with TTL."""
        redis_client = get_redis_client()
        cache_key = f"user_full_data:{student_id}"
        # Convert UUIDs to strings for JSON serialization
        serializable_data = _convert_uuids_to_strings(data)
        await redis_client.setex(
            cache_key, self._settings.REDIS_CACHE_TTL, json.dumps(serializable_data, cls=DateTimeEncoder)
        )
        logger.debug("Cached full data for student %s with TTL %d", student_id, self._settings.REDIS_CACHE_TTL)

    @traced("stats.invalidate_cache", record_args=True, record_result=True)
    async def invalidate_cache(self, student_id: UUID) -> None:
        """Invalidate Redis cache for a student."""
        redis_client = get_redis_client()
        cache_key = f"user_stats:{student_id}"
        await redis_client.delete(cache_key)
        logger.debug("Invalidated cache for student %s", student_id)

    @traced("stats.get_from_db", record_args=True, record_result=True)
    async def get_from_db(self, student_id: UUID) -> dict[str, Any] | None:
        """Get aggregated stats from DB table."""
        result = await fetch_one(q.STUDENT_STATS_AGGREGATED_SELECT, {"student_id": str(student_id)})
        if result and result.get("stats_json"):
            stats = result["stats_json"] if isinstance(result["stats_json"], dict) else json.loads(result["stats_json"])
            # Ensure last_attempt_at is a string for JSON serialization
            if "last_attempt_at" in stats and isinstance(stats["last_attempt_at"], date):
                stats["last_attempt_at"] = stats["last_attempt_at"].isoformat()
            return stats
        return None

    @traced("stats.get_user_stats", record_args=True, record_result=True)
    async def get_user_stats(self, student_id: UUID) -> dict[str, Any]:
        """Get user stats: Redis -> DB -> calculate.

        Args:
            student_id: Student UUID

        Returns:
            Aggregated statistics dictionary
        """
        logger.debug("Getting stats for student %s", student_id)
        # 1. Try Redis
        stats = await self.get_from_redis(student_id)
        if stats:
            return stats
        logger.debug("Stats after Redis: %s", stats)

        logger.debug("Cache miss for student %s, checking DB", student_id)
        # 2. Try DB aggregated table
        stats = await self.get_from_db(student_id)
        if stats:
            await self.save_to_redis(student_id, stats)
            return stats
        logger.debug("Stats after DB: %s", stats)
        logger.debug("DB miss for student %s, calculating stats", student_id)
        # 3. Calculate, save to DB and Redis
        logger.debug("Calculating stats for student %s", student_id)
        logger.debug("Stats before calculation: %s", stats)
        stats = await self.calculate_and_save_aggregated(student_id)
        logger.debug("Stats after calculation: %s", stats)
        await self.save_to_redis(student_id, stats)
        return stats

    @traced("stats.get_student_attempts", record_args=True, record_result=True)
    async def get_student_attempts(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get recent attempts for a student.

        Args:
            student_id: Student UUID

        Returns:
            List of attempt dictionaries
        """
        from app.database import fetch_all
        from app.db import queries as q

        return await fetch_all(q.STUDENT_ATTEMPTS_WITH_CERTIFICATES, {"student_id": student_id})

    @traced("stats.get_raw_user_stats", record_args=True, record_result=True)
    async def get_raw_user_stats(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get raw user statistics data.

        Args:
            student_id: Student UUID

        Returns:
            List of raw user stats dictionaries
        """
        from app.database import fetch_all
        from app.db import queries as q

        return await fetch_all(q.RAW_USER_STATS_BY_STUDENT, {"student_id": student_id})

    @traced("stats.get_raw_attempts", record_args=True, record_result=True)
    async def get_raw_attempts(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get raw attempts data.

        Args:
            student_id: Student UUID

        Returns:
            List of raw attempts dictionaries
        """
        from app.database import fetch_all
        from app.db import queries as q

        return await fetch_all(q.RAW_ATTEMPTS_BY_STUDENT, {"student_id": student_id})


stats_repository = StatsRepository()
