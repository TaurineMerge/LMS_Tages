"""Stats service for user statistics operations."""

import logging
from typing import Any
from uuid import UUID

from app.repositories.stats import stats_repository
from app.services.certificate import certificate_service
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
        """Get aggregated statistics for a user including certificates.

        Args:
            student_id: UUID of the student

        Returns:
            Dictionary containing user statistics and certificates grouped by course
        """

        logger.debug("Getting stats and certificates for student %s", student_id)

        # Try to get from Redis cache first
        cached_data = await self.repo.get_cached_full_data(student_id)
        if cached_data:
            logger.debug("Returning cached data for student %s", student_id)
            return cached_data

        # Get statistics
        stats = await self.repo.get_user_stats(student_id)

        # Get certificates grouped by course
        certificates = await certificate_service.get_certificates_by_student_grouped(student_id)

        # Get recent attempts
        attempts = await self.repo.get_student_attempts(student_id)

        # Get raw data for charts
        raw_user_stats = await self.repo.get_raw_user_stats(student_id)
        raw_attempts = await self.repo.get_raw_attempts(student_id)

        # Combine results
        result = {
            "statistics": stats,
            "certificates": certificates,
            "attempts": attempts,
            "raw_user_stats": raw_user_stats,
            "raw_attempts": raw_attempts,
        }

        # Cache the complete result
        await self.repo.save_full_data_to_redis(student_id, result)

        return result

    @traced("stats_service.refresh_user_statistics", record_args=True, record_result=True)
    async def refresh_user_statistics(self, student_id: UUID) -> dict[str, Any]:
        """Force refresh of user statistics and return with certificates.

        Recalculates from business table, updates DB and Redis cache.

        Args:
            student_id: UUID of the student

        Returns:
            Fresh aggregated statistics with certificates
        """
        stats = await self.repo.calculate_and_save_aggregated(student_id)

        # Get certificates grouped by course
        certificates = await certificate_service.get_certificates_by_student_grouped(student_id)

        # Get recent attempts
        attempts = await self.repo.get_student_attempts(student_id)

        # Get raw data for charts
        raw_user_stats = await self.repo.get_raw_user_stats(student_id)
        raw_attempts = await self.repo.get_raw_attempts(student_id)

        # Combine results
        result = {
            "statistics": stats,
            "certificates": certificates,
            "attempts": attempts,
            "raw_user_stats": raw_user_stats,
            "raw_attempts": raw_attempts,
        }

        # Cache the complete result
        await self.repo.save_full_data_to_redis(student_id, result)

        return result
