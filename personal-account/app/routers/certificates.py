"""Certificate API endpoints."""

from uuid import UUID

from fastapi import APIRouter, Depends, Query, status

from app.core.security import TokenPayload, get_current_user, require_roles
from app.schemas.certificate import certificate_create, certificate_response
from app.services.certificate import certificate_service
from app.telemetry import traced

router = APIRouter(prefix="/certificates", tags=["Certificates"])


@router.get(
    "",
    response_model=list[certificate_response],
    summary="Получить список сертификатов",
    description="Возвращает список сертификатов с опциональной фильтрацией по студенту и курсу",
)
@traced("router.certificates.get_certificates")
async def get_certificates(
    student_id: UUID | None = Query(default=None, description="Фильтр по студенту"),
    course_id: UUID | None = Query(default=None, description="Фильтр по курсу"),
    current_user: TokenPayload = Depends(get_current_user),
):
    """Get list of certificates with optional filters."""
    return await certificate_service.get_certificates(student_id, course_id)


@router.get(
    "/{certificate_id}",
    response_model=certificate_response,
    summary="Получить сертификат по ID",
    description="Возвращает данные сертификата по указанному ID",
)
@traced("router.certificates.get_certificate")
async def get_certificate(certificate_id: UUID, current_user: TokenPayload = Depends(get_current_user)):
    """Get certificate by ID."""
    return await certificate_service.get_certificate(certificate_id)


@router.post(
    "",
    response_model=certificate_response,
    status_code=status.HTTP_201_CREATED,
    summary="Создать новый сертификат",
    description="Создает новый сертификат для студента",
)
@traced("router.certificates.create_certificate")
async def create_certificate(
    certificate: certificate_create, current_user: TokenPayload = Depends(require_roles("admin"))
):
    """Create a new certificate."""
    return await certificate_service.create_certificate(certificate)


@router.delete(
    "/{certificate_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Удалить сертификат",
    description="Удаляет сертификат из системы",
)
@traced("router.certificates.delete_certificate")
async def delete_certificate(certificate_id: UUID):
    """Delete certificate by ID."""
    await certificate_service.delete_certificate(certificate_id)
    return None
