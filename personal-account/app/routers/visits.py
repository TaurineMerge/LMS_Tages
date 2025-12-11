"""Visit API endpoints."""

from uuid import UUID

from fastapi import APIRouter, Depends, Query, status

from app.core.security import TokenPayload, get_current_user, require_roles
from app.schemas.visit import visit_create, visit_response
from app.services.visit import visit_service
from app.telemetry import traced

router = APIRouter(prefix="/visits", tags=["Visits"])


@router.get(
    "",
    response_model=list[visit_response],
    summary="Получить список посещений уроков",
    description="Возвращает список посещений с опциональной фильтрацией по студенту и уроку",
)
@traced("router.visits.get_visits", record_args=True, record_result=True)
async def get_visits(
    student_id: UUID | None = Query(default=None, description="Фильтр по студенту"),
    lesson_id: UUID | None = Query(default=None, description="Фильтр по уроку"),
    current_user: TokenPayload = Depends(get_current_user),
):
    """Get list of visits with optional filters."""
    return await visit_service.get_visits(student_id, lesson_id)


@router.get(
    "/{visit_id}",
    response_model=visit_response,
    summary="Получить посещение по ID",
    description="Возвращает данные посещения по указанному ID",
)
@traced("router.visits.get_visit", record_args=True, record_result=True)
async def get_visit(visit_id: UUID, current_user: TokenPayload = Depends(get_current_user)):
    """Get visit by ID."""
    return await visit_service.get_visit(visit_id)


@router.post(
    "",
    response_model=visit_response,
    status_code=status.HTTP_201_CREATED,
    summary="Зарегистрировать посещение урока",
    description="Регистрирует посещение урока студентом",
)
@traced("router.visits.create_visit", record_args=True, record_result=True)
async def create_visit(visit: visit_create, current_user: TokenPayload = Depends(require_roles("admin", "teacher"))):
    """Create a new visit record."""
    return await visit_service.create_visit(visit)


@router.delete(
    "/{visit_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Удалить запись о посещении",
    description="Удаляет запись о посещении урока",
)
@traced("router.visits.delete_visit", record_args=True, record_result=True)
async def delete_visit(visit_id: UUID):
    """Delete visit by ID."""
    await visit_service.delete_visit(visit_id)
    return None
