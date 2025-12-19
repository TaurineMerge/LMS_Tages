"""Stats service for user statistics operations."""

import logging
from typing import Any
from uuid import UUID

from app.repositories.stats import stats_repository
from app.telemetry import traced

logger = logging.getLogger(__name__)


class StatsService:
    """Service for user statistics operations.

    Provides high-level operations for retrieving user statistics
    with caching support.
    """

    def __init__(self):
        self.repo = stats_repository

    @traced("stats_service.get_user_statistics", record_args=True, record_result=True)
    async def get_user_statistics(self, student_id: UUID) -> dict[str, Any]:
        """Get aggregated statistics for a user.

        Retrieves statistics from cache if available, otherwise calculates
        from database and caches the result.

        Args:
            student_id: UUID of the student

        Returns:
            Dictionary containing user statistics
        """
        return await self.repo.get_user_stats(student_id)
