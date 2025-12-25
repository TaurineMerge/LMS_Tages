"""Certificate schemas."""

from datetime import date
from uuid import UUID

from pydantic import BaseModel, ConfigDict, FieldValidationInfo, HttpUrl, field_validator

from app.schemas.validators import ensure_safe_string


class certificate_base(BaseModel):
    """Base certificate schema."""

    content: str | None = None
    student_id: UUID
    course_id: UUID

    @field_validator("content", mode="before")
    @classmethod
    def _validate_content(cls, value: str | None, info: FieldValidationInfo) -> str | None:
        return ensure_safe_string(value, info.field_name)


class certificate_create(certificate_base):
    """Schema for creating a certificate."""

    test_attempt_id: UUID | None = None


class certificate_response(BaseModel):
    """Schema for certificate response."""

    model_config = ConfigDict(from_attributes=True)

    id: UUID
    certificate_number: int
    created_at: date
    content: str | None = None
    student_id: UUID
    course_id: UUID
    test_attempt_id: UUID | None = None


class certificate_download(BaseModel):
    """Schema for certificate download response."""

    model_config = ConfigDict(from_attributes=True)

    id: UUID
    certificate_number: int
    student_id: UUID
    course_id: UUID
    download_url: HttpUrl | None = None
    created_at: date


class certificate_list_item(BaseModel):
    """Schema for certificate list item."""

    model_config = ConfigDict(from_attributes=True)

    id: UUID
    certificate_number: int
    created_at: date
    course_id: UUID
    student_id: UUID
