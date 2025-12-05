"""Certificate API endpoints."""
from uuid import UUID

from fastapi import APIRouter, Query, status

from app.schemas.certificate import CertificateCreate, CertificateResponse
from app.services.certificate import certificate_service

router = APIRouter(prefix="/certificates", tags=["Certificates"])


@router.get(
    "",
    response_model=list[CertificateResponse],
    summary="Получить список сертификатов",
    description="Возвращает список сертификатов с опциональной фильтрацией по студенту и курсу"
)
async def get_certificates(
    student_id: UUID | None = Query(default=None, description="Фильтр по студенту"),
    course_id: UUID | None = Query(default=None, description="Фильтр по курсу")
):
    """Get list of certificates with optional filters."""
    return await certificate_service.get_certificates(student_id, course_id)


@router.get(
    "/{certificate_id}",
    response_model=CertificateResponse,
    summary="Получить сертификат по ID",
    description="Возвращает данные сертификата по указанному ID"
)
async def get_certificate(certificate_id: UUID):
    """Get certificate by ID."""
    return await certificate_service.get_certificate(certificate_id)


@router.post(
    "",
    response_model=CertificateResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Создать новый сертификат",
    description="Создает новый сертификат для студента"
)
async def create_certificate(certificate: CertificateCreate):
    """Create a new certificate."""
    return await certificate_service.create_certificate(certificate)


@router.delete(
    "/{certificate_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Удалить сертификат",
    description="Удаляет сертификат из системы"
)
async def delete_certificate(certificate_id: UUID):
    """Delete certificate by ID."""
    await certificate_service.delete_certificate(certificate_id)
    return None
