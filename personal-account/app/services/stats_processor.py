"""StatsProcessor: validates and processes raw payloads into business tables."""

import logging
from uuid import UUID

from app.clients.validation.contract_manager import ContractManager
from app.repositories.certificate import certificate_repository
from app.repositories.stats import stats_repository
from app.repositories.student import student_repository
from app.services.cert_service import CertificateGenerationError, get_certificate_service
from app.telemetry import traced

logger = logging.getLogger(__name__)


class StatsProcessor:
    """Processor that validates raw payloads and writes business records.

    Orchestrates validation (ContractManager), lookups (student repository)
    and persistence (StatsRepository). After processing, recalculates
    aggregated stats and updates cache.
    """

    def __init__(self, db=None):
        self.certificate_repo = certificate_repository
        self.repo = stats_repository
        self.cm = ContractManager()
        self.student_repo = student_repository
        self.certificate_service = get_certificate_service()

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

    @traced("stats_processor.check_and_generate_certificates_for_student", record_args=True, record_result=True)
    async def check_and_generate_certificates_for_student(self, student_id: UUID) -> dict[str, int]:
        """Check for passing attempts without certificates for a specific student and generate them."""
        stats = {"generated": 0, "failed": 0}

        passing_attempts = await self.certificate_repo.get_passing_attempts_without_certificates_for_student(student_id)

        for attempt in passing_attempts:
            try:
                # Generate certificate using cert_service
                certificate_data = await self.certificate_service.generate_certificate(
                    student_id=attempt["student_id"],
                    course_id=attempt["course_id"],
                    test_attempt_id=attempt["id"],
                    score=attempt["score"],
                    max_score=attempt["max_score"],
                    course_name=attempt["course_name"],
                )

                # Save certificate record
                await self.certificate_repo.create_certificate(certificate_data)

                stats["generated"] += 1
                logger.info("Generated certificate for attempt %s", attempt["id"])

            except Exception as e:
                logger.error("Failed to generate certificate for attempt %s: %s", attempt["id"], e)
                stats["failed"] += 1

        return stats

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

    @traced("stats_processor.check_and_generate_certificates", record_args=True, record_result=True)
    async def check_and_generate_certificates(self) -> dict[str, int]:
        """Check for passing attempts without certificates and generate them."""
        stats = {"generated": 0, "failed": 0}

        passing_attempts = await self.certificate_repo.get_passing_attempts_without_certificates()

        for attempt in passing_attempts:
            try:
                # Generate certificate using cert_service
                certificate_data = await self.certificate_service.generate_certificate(
                    student_id=attempt["student_id"],
                    course_id=attempt["course_id"],
                    test_attempt_id=attempt["id"],
                    score=attempt["score"],
                    max_score=attempt["max_score"],
                    course_name=attempt["course_name"],
                )

                # Save certificate record
                await self.certificate_repo.create_certificate(certificate_data)

                stats["generated"] += 1
                logger.info("Generated certificate for attempt %s", attempt["id"])

            except Exception as e:
                logger.error("Failed to generate certificate for attempt %s: %s", attempt["id"], e)
                stats["failed"] += 1

        return stats


stats_processor = StatsProcessor()
