"""Certificate schemas."""
from datetime import date
from uuid import UUID

from pydantic import BaseModel, ConfigDict


class CertificateBase(BaseModel):
    """Base certificate schema."""
    
    content: str | None = None
    student_id: UUID
    course_id: UUID


class CertificateCreate(CertificateBase):
    """Schema for creating a certificate."""
    
    test_attempt_id: UUID | None = None


class CertificateResponse(BaseModel):
    """Schema for certificate response."""
    
    model_config = ConfigDict(from_attributes=True)
    
    id: UUID
    certificate_number: int
    created_at: date
    content: str | None = None
    student_id: UUID
    course_id: UUID
    test_attempt_id: UUID | None = None
