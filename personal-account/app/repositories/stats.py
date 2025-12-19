"""Business repository for stats-related DB operations.

This module contains the database operations that populate the
business table `tests.test_attempt_b` from raw payloads and manage
aggregated statistics with Redis caching.
"""

import json
import logging
from typing import Any
from uuid import UUID

from app.config import get_settings
from app.database import execute, execute_returning, fetch_all, fetch_one
from app.db import queries as q
from app.redis_client import get_redis_client
from app.telemetry import traced

logger = logging.getLogger(__name__)


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
        params = {
            "id": str(UUID(payload["attempt_id"])),
            "student_id": str(UUID(payload["student_id"])),
            "test_id": str(UUID(payload["test_id"])) if payload.get("test_id") else None,
            "date_of_attempt": payload.get("date_of_attempt"),
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

        # Save to aggregated table
        params = {
            "student_id": str(student_id),
            "total_attempts": stats["total_attempts"],
            "passed_attempts": stats["passed_attempts"],
            "failed_attempts": stats["failed_attempts"],
            "avg_score": stats["avg_score"],
            "total_tests_taken": stats["total_tests_taken"],
            "last_attempt_at": stats["last_attempt_at"],
            "stats_json": json.dumps(stats),
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
            logger.debug("Cache hit for student %s", student_id)
            return json.loads(cached)
        return None

    @traced("stats.save_to_redis", record_args=True, record_result=True)
    async def save_to_redis(self, student_id: UUID, stats: dict[str, Any]) -> None:
        """Save stats to Redis with TTL."""
        redis_client = get_redis_client()
        cache_key = f"user_stats:{student_id}"
        await redis_client.setex(cache_key, self._settings.REDIS_CACHE_TTL, json.dumps(stats))
        logger.debug("Cached stats for student %s with TTL %d", student_id, self._settings.REDIS_CACHE_TTL)

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
            return result["stats_json"] if isinstance(result["stats_json"], dict) else json.loads(result["stats_json"])
        return None

    @traced("stats.get_user_stats", record_args=True, record_result=True)
    async def get_user_stats(self, student_id: UUID) -> dict[str, Any]:
        """Get user stats: Redis -> DB -> calculate.

        Args:
            student_id: Student UUID

        Returns:
            Aggregated statistics dictionary
        """
        # 1. Try Redis
        stats = await self.get_from_redis(student_id)
        if stats:
            return stats

        # 2. Try DB aggregated table
        stats = await self.get_from_db(student_id)
        if stats:
            await self.save_to_redis(student_id, stats)
            return stats

        # 3. Calculate, save to DB and Redis
        stats = await self.calculate_and_save_aggregated(student_id)
        await self.save_to_redis(student_id, stats)
        return stats


stats_repository = StatsRepository()
