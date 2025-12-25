"""StatsWorker service for background processing of statistics data.

This service runs periodic jobs to fetch data from the testing service
and process raw payloads into business tables using AsyncIOScheduler.
"""

import logging
from uuid import UUID

from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.interval import IntervalTrigger

from app.clients.testing_client import TestingClient
from app.config import get_settings
from app.repositories.integration import integration_repository
from app.repositories.student import student_repository
from app.services.stats_processor import StatsProcessor
from app.telemetry import traced

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)


class StatsWorker:
    """Background worker for statistics integration and processing.

    Runs periodic jobs:
    - fetch_from_testing: fetches new data from testing client and saves to integration.raw_*
    - process_raws: processes raw attempts and user stats into business tables

    Uses AsyncIOScheduler for job scheduling. Intervals are configurable via settings.
    """

    def __init__(
        self,
        testing_client: TestingClient | None = None,
        processor: StatsProcessor | None = None,
    ):
        self.testing_client = testing_client or TestingClient(get_settings().TESTING_BASE_URL)
        self.processor = processor or StatsProcessor()
        self.scheduler = AsyncIOScheduler()
        settings = get_settings()
        self.fetch_interval = settings.STATS_WORKER_FETCH_INTERVAL
        self.process_interval = settings.STATS_WORKER_PROCESS_INTERVAL

    @traced("stats_worker.start", record_args=True, record_result=True)
    def start(self) -> None:
        """Start the worker scheduler with periodic jobs."""
        self.scheduler.add_job(
            self._wrap(self.fetch_from_testing),
            IntervalTrigger(seconds=self.fetch_interval),
            id="fetch_testing",
            replace_existing=True,
        )
        self.scheduler.add_job(
            self._wrap(self.process_raws),
            IntervalTrigger(seconds=self.process_interval),
            id="process_raws",
            replace_existing=True,
        )
        self.scheduler.add_job(
            self._wrap(self._generate_certificates),
            IntervalTrigger(seconds=60),  # Check for certificates every minute
            id="generate_certificates",
            replace_existing=True,
        )
        if not self.scheduler.running:
            self.scheduler.start()
            logger.info(
                "StatsWorker started with fetch every %ds, process every %ds",
                self.fetch_interval,
                self.process_interval,
            )

    @traced("stats_worker.stop", record_args=True, record_result=True)
    def stop(self) -> None:
        """Stop the worker scheduler."""
        if self.scheduler.running:
            self.scheduler.shutdown()
            logger.info("StatsWorker stopped")

    def _wrap(self, coro):
        """Wrap a coroutine to handle exceptions in scheduled jobs."""

        async def job_wrapper():
            try:
                await coro()
            except Exception:
                logger.exception("Worker job failed")

        return job_wrapper

    @traced("stats_worker.fetch_for_student", record_args=True, record_result=True)
    async def fetch_for_student(self, student_id: UUID) -> None:
        """Fetch data for a specific student."""
        logger.info("Fetching data for student %s", student_id)
        try:
            # Fetch stats
            stats_payload = await self.testing_client.get_user_stats(student_id)
            await integration_repository.save_raw_user_stats(student_id, stats_payload)

            # Fetch attempts
            attempts_payload = await self.testing_client.get_user_attempts(student_id)
            await integration_repository.save_raw_attempts(student_id, attempts_payload)

            logger.info("Successfully fetched data for student %s", student_id)
        except Exception as e:
            logger.error("Failed to fetch data for student %s: %s", student_id, e)
            raise

    @traced("stats_worker.fetch_from_testing", record_args=True, record_result=True)
    async def fetch_from_testing(self) -> None:
        """Fetch new data from testing service and save to integration.raw_*.

        Iterates over all students, fetches their stats and attempts,
        and persists raw payloads. In production, use paging and checkpoints
        to avoid fetching all users every time.
        """
        logger.info("Starting fetch from testing service")
        # Get all student IDs (simplified; in production, use paging/checkpoints)
        students = await student_repository.get_paginated(page=1, limit=3)  # Adjust limit as needed
        student_ids = [s["id"] for s in students[0]]

        for student_id in student_ids:
            try:
                # Fetch user stats
                stats_payload = await self.testing_client.fetch_user_stats(student_id)
                await self.testing_client.repo.insert_raw_user_stats(
                    student_id, stats_payload
                )  # Assuming repo is accessible

                # Fetch user attempts
                attempts_list = await self.testing_client.fetch_user_attempts(student_id)
                for attempt in attempts_list:
                    attempt_id = attempt["attempt_id"]
                    # Fetch attempt detail
                    attempt_detail = await self.testing_client.fetch_attempt_detail(attempt_id)
                    await self.testing_client.repo.insert_raw_attempt(attempt_id, student_id, None, attempt_detail)

            except Exception:
                logger.exception("Failed fetching for student %s", student_id)

        logger.info("Completed fetch from testing service")

    @traced("stats_worker.process_raws", record_args=True, record_result=True)
    async def process_raws(self) -> None:
        """Process raw payloads into business tables.

        Processes user stats first, then attempts.
        """
        logger.info("Starting raw processing")
        # Process user stats
        user_stats_result = await self.processor.process_raw_user_stats()
        logger.info("Processed user stats: %s", user_stats_result)

        # Process attempts
        attempts_result = await self.processor.process_raw_attempts()
        logger.info("Processed attempts: %s", attempts_result)
        logger.info("Completed raw processing")

    @traced("stats_worker._generate_certificates", record_args=True, record_result=True)
    async def _generate_certificates(self) -> None:
        """Generate certificates for successful test attempts.

        Runs periodically to check for passing attempts without certificates
        and generate certificates for them.
        """
        logger.info("Starting certificate generation check")
        try:
            result = await self.processor.check_and_generate_certificates()
            logger.info(
                "Certificate generation completed: generated=%d, failed=%d",
                result.get("certificates_generated", 0),
                result.get("failed", 0),
            )
        except Exception:
            logger.exception("Failed to check and generate certificates")
