"""Student API endpoints."""

from uuid import UUID

from fastapi import APIRouter, Depends, Query, status

from app.core.security import TokenPayload, get_current_user, require_roles
from app.schemas.common import paginated_response
from app.schemas.student import student_create, student_response, student_update
from app.services.student import student_service
from app.telemetry import traced

router = APIRouter(prefix="/students", tags=["Students"])


@router.get(
    "",
    response_model=paginated_response[student_response],
    summary="Получить список студентов",
    description="Возвращает пагинированный список всех студентов",
)
@traced("router.students.get_students")
async def get_students(
    page: int = Query(default=1, ge=1, description="Номер страницы"),
    limit: int = Query(default=20, ge=1, le=100, description="Количество элементов на странице"),
    current_user: TokenPayload = Depends(get_current_user),
):
    """Get paginated list of students."""
    return await student_service.get_students(page, limit)


@router.get(
    "/{student_id}",
    response_model=student_response,
    summary="Получить студента по ID",
    description="Возвращает данные студента по указанному ID",
)
@traced("router.students.get_student")
async def get_student(student_id: UUID, current_user: TokenPayload = Depends(get_current_user)):
    """Get student by ID."""
    return await student_service.get_student(student_id)


@router.post(
    "",
    response_model=student_response,
    status_code=status.HTTP_201_CREATED,
    summary="Создать нового студента",
    description="Создает нового студента в системе",
)
@traced("router.students.create_student")
async def create_student(student: student_create, current_user: TokenPayload = Depends(require_roles("admin"))):
    """Create a new student."""
    return await student_service.create_student(student)


@router.put(
    "/{student_id}",
    response_model=student_response,
    summary="Обновить данные студента",
    description="Обновляет данные существующего студента",
)
@traced("router.students.update_student")
async def update_student(student_id: UUID, student: student_update):
    """Update student by ID."""
    return await student_service.update_student(student_id, student)


@router.delete(
    "/{student_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Удалить студента",
    description="Удаляет студента из системы",
)
@traced("router.students.delete_student")
async def delete_student(student_id: UUID):
    """Delete student by ID."""
    await student_service.delete_student(student_id)
    return None
