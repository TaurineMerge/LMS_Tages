"""Stats service for user statistics operations."""

import logging
from typing import Any
from uuid import UUID

from app.repositories.stats import stats_repository
from app.telemetry import traced

logger = logging.getLogger(__name__)


class StatsService:
    """Service for user statistics with Redis caching.

    Provides high-level operations for retrieving user statistics.
    Flow: Redis cache -> DB aggregated table -> calculate from business table.
    """

    def __init__(self):
        self.repo = stats_repository

    @traced("stats_service.get_user_statistics", record_args=True, record_result=True)
    async def get_user_statistics(self, student_id: UUID) -> dict[str, Any]:
        """Get aggregated statistics for a user.

        Args:
            student_id: UUID of the student

        Returns:
            Dictionary containing user statistics
        """
        return await self.repo.get_user_stats(student_id)

    @traced("stats_service.refresh_user_statistics", record_args=True, record_result=True)
    async def refresh_user_statistics(self, student_id: UUID) -> dict[str, Any]:
        """Force refresh of user statistics.

        Recalculates from business table, updates DB and Redis cache.

        Args:
            student_id: UUID of the student

        Returns:
            Fresh aggregated statistics
        """
        stats = await self.repo.calculate_and_save_aggregated(student_id)
        await self.repo.save_to_redis(student_id, stats)
        return stats


stats_service = StatsService()
