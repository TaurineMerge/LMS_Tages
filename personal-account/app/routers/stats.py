"""Statistics API router."""

import logging
from uuid import UUID

from fastapi import APIRouter, HTTPException

from app.services import stats_service, stats_worker
from app.telemetry import traced

router = APIRouter(prefix="/stats", tags=["Statistics"])

logger = logging.getLogger(__name__)


@router.get("/{student_id}")
@traced("api.get_student_stats", record_args=True, record_result=True)
async def get_student_stats(student_id: UUID):
    """Get aggregated statistics for a student."""
    try:
        logger.info("Fetching stats for student %s", student_id)
        stats = await stats_service.get_user_statistics(student_id)
        return {"status": "ok", "data": stats}
    except Exception as e:
        logger.error("Failed to get stats for %s: %s", student_id, e)
        raise HTTPException(status_code=500, detail=str(e)) from e


@router.post("/{student_id}/refresh")
@traced("api.refresh_student_stats", record_args=True, record_result=True)
async def refresh_student_stats(student_id: UUID):
    """Force refresh statistics for a student."""
    try:
        logger.info("Refreshing stats for student %s", student_id)
        stats = await stats_service.refresh_user_statistics(student_id)
        return {"status": "ok", "data": stats}
    except Exception as e:
        logger.error("Failed to refresh stats for %s: %s", student_id, e)
        raise HTTPException(status_code=500, detail=str(e)) from e


@router.post("/worker/run-fetch")
@traced("api.worker_run_fetch")
async def run_worker_fetch():
    """Manually trigger worker fetch from testing service."""
    try:
        logger.info("Manual trigger: worker fetch")
        await stats_worker.fetch_from_testing()
        return {"status": "ok", "message": "Fetch completed"}
    except Exception as e:
        logger.error("Manual fetch failed: %s", e)
        raise HTTPException(status_code=500, detail=str(e)) from e


@router.post("/worker/run-process")
@traced("api.worker_run_process")
async def run_worker_process():
    """Manually trigger worker processing of raw data."""
    try:
        logger.info("Manual trigger: worker process")
        await stats_worker.process_raws()
        return {"status": "ok", "message": "Processing completed"}
    except Exception as e:
        logger.error("Manual processing failed: %s", e)
        raise HTTPException(status_code=500, detail=str(e)) from e
