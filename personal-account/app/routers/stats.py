"""Statistics API router."""

import logging as logger
from uuid import UUID

from fastapi import APIRouter, HTTPException

from app.services.stats_service import stats_service
from app.telemetry import traced

router = APIRouter(prefix="/stats", tags=["Statistics"])


@router.get("/{student_id}")
@traced("api.get_student_stats", record_args=True, record_result=True)
async def get_student_stats(student_id: UUID):
    """Get aggregated statistics for a student.

    Args:
        student_id: UUID of the student

    Returns:
        Aggregated statistics dictionary
    """
    try:
        logger.log("Fetching stats for student %s", student_id)
        stats = await stats_service.get_user_statistics(student_id)
        logger.log("Fetched stats for student %s: %s", student_id, stats)
        return {"status": "ok", "data": stats}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) from e


@router.post("/{student_id}/refresh")
@traced("api.refresh_student_stats", record_args=True, record_result=True)
async def refresh_student_stats(student_id: UUID):
    """Force refresh statistics for a student.

    Args:
        student_id: UUID of the student

    Returns:
        Fresh aggregated statistics
    """
    try:
        logger.log("Refreshing stats for student %s", student_id)
        stats = await stats_service.refresh_user_statistics(student_id)
        logger.log("Refreshed stats for student %s: %s", student_id, stats)
        return {"status": "ok", "data": stats}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e)) from e
