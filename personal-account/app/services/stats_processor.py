"""StatsProcessor: validates and processes raw payloads into business tables."""

import logging
from uuid import UUID

from app.clients.validation.contract_manager import ContractManager
from app.repositories.stats import stats_repository
from app.repositories.student import student_repository
from app.telemetry import traced

logger = logging.getLogger(__name__)


class StatsProcessor:
    """Processor that validates raw payloads and writes business records.

    Orchestrates validation (ContractManager), lookups (student repository)
    and persistence (StatsRepository). After processing, recalculates
    aggregated stats and updates cache.
    """

    def __init__(self, db=None):
        self.repo = stats_repository
        self.cm = ContractManager()
        self.student_repo = student_repository

    @traced("stats_processor.process_raw_attempts", record_args=True, record_result=True)
    async def process_raw_attempts(self, batch_size: int = 100) -> dict[str, int]:
        """Process unprocessed raw attempts.

        Args:
            batch_size: Maximum number of raw rows to fetch and process.

        Returns:
            Dict with counters: processed and failed.
        """
        stats = {"processed": 0, "failed": 0}
        affected_students: set[UUID] = set()
        raws = await self.repo.get_unprocessed_attempts(limit=batch_size)

        for raw in raws:
            raw_id = raw["id"]
            payload = raw["payload"]
            try:
                await self.cm.validate_attempt_detail(payload)
                student_id = UUID(payload.get("student_id"))
                student = await self.student_repo.get_by_id(student_id)
                if not student:
                    raise ValueError(f"Student {student_id} not found")
                await self.repo.upsert_attempt(payload)
                await self.repo.mark_attempt_processed(raw_id, None)
                affected_students.add(student_id)
                stats["processed"] += 1
            except Exception as e:
                logger.exception("Failed to process raw attempt %s", raw_id)
                await self.repo.mark_attempt_processed(raw_id, str(e))
                stats["failed"] += 1

        # Recalculate and cache for affected students
        for student_id in affected_students:
            try:
                aggregated = await self.repo.calculate_and_save_aggregated(student_id)
                await self.repo.save_to_redis(student_id, aggregated)
            except Exception:
                logger.exception("Failed to recalculate stats for student %s", student_id)

        return stats

    @traced("stats_processor.process_raw_user_stats", record_args=True, record_result=True)
    async def process_raw_user_stats(self, batch_size: int = 100) -> dict[str, int]:
        """Process unprocessed user-level raw stats.

        Args:
            batch_size: Maximum number of raw rows to fetch and process.

        Returns:
            Dict with counters: processed and failed.
        """
        stats = {"processed": 0, "failed": 0}
        raws = await self.repo.get_unprocessed_user_stats(limit=batch_size)

        for raw in raws:
            raw_id = raw["id"]
            payload = raw["payload"]
            try:
                await self.cm.validate_user_stats(payload)
                await self.repo.mark_user_stats_processed(raw_id, None)
                stats["processed"] += 1
            except Exception as e:
                logger.exception("Failed to process user stats %s", raw_id)
                await self.repo.mark_user_stats_processed(raw_id, str(e))
                stats["failed"] += 1

        return stats
