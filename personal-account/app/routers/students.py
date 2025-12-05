"""Student API endpoints."""
from uuid import UUID

from fastapi import APIRouter, Query, status

from app.schemas.student import StudentCreate, StudentUpdate, StudentResponse
from app.schemas.common import PaginatedResponse
from app.services.student import student_service

router = APIRouter(prefix="/students", tags=["Students"])


@router.get(
    "",
    response_model=PaginatedResponse[StudentResponse],
    summary="Получить список студентов",
    description="Возвращает пагинированный список всех студентов"
)
async def get_students(
    page: int = Query(default=1, ge=1, description="Номер страницы"),
    limit: int = Query(default=20, ge=1, le=100, description="Количество элементов на странице")
):
    """Get paginated list of students."""
    return await student_service.get_students(page, limit)


@router.get(
    "/{student_id}",
    response_model=StudentResponse,
    summary="Получить студента по ID",
    description="Возвращает данные студента по указанному ID"
)
async def get_student(student_id: UUID):
    """Get student by ID."""
    return await student_service.get_student(student_id)


@router.post(
    "",
    response_model=StudentResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Создать нового студента",
    description="Создает нового студента в системе"
)
async def create_student(student: StudentCreate):
    """Create a new student."""
    return await student_service.create_student(student)


@router.put(
    "/{student_id}",
    response_model=StudentResponse,
    summary="Обновить данные студента",
    description="Обновляет данные существующего студента"
)
async def update_student(student_id: UUID, student: StudentUpdate):
    """Update student by ID."""
    return await student_service.update_student(student_id, student)


@router.delete(
    "/{student_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Удалить студента",
    description="Удаляет студента из системы"
)
async def delete_student(student_id: UUID):
    """Delete student by ID."""
    await student_service.delete_student(student_id)
    return None
