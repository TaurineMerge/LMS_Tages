# Service layer

from .certificate import certificate_service
from .stats_service import StatsService
from .stats_worker import StatsWorker
from .student import student_service
from .visit import visit_service

stats_service = StatsService()
stats_worker = StatsWorker()

__all__ = [
    "certificate_service",
    "stats_service",
    "stats_worker",
    "student_service",
    "visit_service",
]
