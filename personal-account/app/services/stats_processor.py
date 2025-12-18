"""
StatsProcessor:
- валидирует payload через ContractManager
- трансформирует и записывает в бизнес таблицу tests.test_attempt_b
- отмечает raw записи как processed
"""

import logging
from uuid import UUID

from app.clients.validation.contract_manager import ContractManager
from app.repositories.stats import stats_repository
from app.repositories.student import student_repository  # предполагается что есть
from app.telemetry import traced

logger = logging.getLogger(__name__)


class StatsProcessor:
    """Processor that validates raw payloads and writes business records.

    The processor orchestrates validation (via ContractManager), business
    lookups (student repository) and persistence (StatsRepository).
    """

    def __init__(self, db=None):
        # db parameter kept for compatibility but not used; repositories are singletons
        self.repo = stats_repository
        self.cm = ContractManager()
        self.student_repo = student_repository

    @traced("stats.process_raw_attempts", record_args=True, record_result=True)
    async def process_raw_attempts(self, batch_size: int = 100) -> dict[str, int]:
        """Process unprocessed raw attempts.

        Parameters
        - batch_size: maximum number of raw rows to fetch and process.

        Returns a dict with counters: processed and failed.
        """
        stats = {"processed": 0, "failed": 0}
        raws = await self.repo.get_unprocessed_attempts(limit=batch_size)
        for raw in raws:
            raw_id = raw["id"]
            payload = raw["payload"]
            try:
                # Валидация через ContractManager (attempt_detail)
                await self.cm.validate_attempt_detail(payload)
                # Проверка студента
                student_id = UUID(payload.get("student_id"))
                student = await self.student_repo.get_by_id(student_id)
                if not student:
                    raise ValueError(f"Student {student_id} not found")
                # Вставка/обновление в tests.test_attempt_b через репозиторий
                await self.repo.upsert_attempt(payload)
                await self.repo.mark_attempt_processed(raw_id, None)
                stats["processed"] += 1
            except Exception as e:
                logger.exception("Failed to process raw attempt %s", raw_id)
                await self.repo.mark_attempt_processed(raw_id, str(e))
                stats["failed"] += 1
        return stats

    @traced("stats.process_raw_user_stats", record_args=True, record_result=True)
    async def process_raw_user_stats(self, batch_size: int = 100) -> dict[str, int]:
        """Process unprocessed user-level aggregated raw stats.

        This method validates the payload and marks the raw record as processed.
        Business-side aggregation/upserts should live in StatsRepository when needed.
        """
        stats = {"processed": 0, "failed": 0}
        raws = await self.repo.get_unprocessed_user_stats(limit=batch_size)
        for raw in raws:
            raw_id = raw["id"]
            payload = raw["payload"]
            try:
                await self.cm.validate_user_stats(payload)
                # Здесь можно агрегировать/обновить дополнительные справочники при необходимости
                await self.repo.mark_user_stats_processed(raw_id, None)
                stats["processed"] += 1
            except Exception as e:
                logger.exception("Failed to process user stats %s", raw_id)
                await self.repo.mark_user_stats_processed(raw_id, str(e))
                stats["failed"] += 1
        return stats
