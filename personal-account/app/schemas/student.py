"""Student schemas."""

from datetime import date, datetime
from typing import Any
from uuid import UUID

from pydantic import (
    BaseModel,
    ConfigDict,
    EmailStr,
    Field,
    FieldValidationInfo,
    field_validator,
)

from app.schemas.validators import ensure_safe_mapping, ensure_safe_string


class student_base(BaseModel):
    """Base student schema with common fields."""

    name: str = Field(..., min_length=1, max_length=100, description="Имя студента")
    surname: str = Field(..., min_length=1, max_length=100, description="Фамилия студента")
    birth_date: date | None = Field(None, description="Дата рождения")
    avatar: str | None = Field(None, max_length=500, description="URL аватара")
    contacts: dict[str, Any] | None = Field(default_factory=dict, description="Контактная информация")
    email: EmailStr = Field(..., description="Email адрес")
    phone: str | None = Field(None, max_length=20, description="Номер телефона")

    @field_validator("name", "surname", "avatar", "phone", mode="before")
    @classmethod
    def _validate_safe_strings(cls, value: str | None, info: FieldValidationInfo) -> str | None:
        return ensure_safe_string(value, info.field_name)

    @field_validator("email", mode="before")
    @classmethod
    def _validate_email_string(cls, value: str | None, info: FieldValidationInfo) -> str | None:
        return ensure_safe_string(value, info.field_name)

    @field_validator("contacts", mode="before")
    @classmethod
    def _validate_contacts(cls, value: dict[str, Any] | None, info: FieldValidationInfo) -> dict[str, Any] | None:
        return ensure_safe_mapping(value, info.field_name)


class student_create(student_base):
    """Schema for creating a student."""

    pass


class student_update(BaseModel):
    """Schema for updating a student. All fields are optional."""

    model_config = ConfigDict(extra="forbid")

    name: str | None = Field(None, min_length=1, max_length=100)
    surname: str | None = Field(None, min_length=1, max_length=100)
    birth_date: date | None = None
    avatar: str | None = Field(None, max_length=500)
    contacts: dict[str, Any] | None = None
    email: EmailStr | None = None
    phone: str | None = Field(None, max_length=20)

    @field_validator("name", "surname", "avatar", "phone", mode="before")
    @classmethod
    def _validate_update_strings(cls, value: str | None, info: FieldValidationInfo) -> str | None:
        return ensure_safe_string(value, info.field_name)

    @field_validator("email", mode="before")
    @classmethod
    def _validate_update_email(cls, value: str | None, info: FieldValidationInfo) -> str | None:
        return ensure_safe_string(value, info.field_name)

    @field_validator("contacts", mode="before")
    @classmethod
    def _validate_update_contacts(
        cls, value: dict[str, Any] | None, info: FieldValidationInfo
    ) -> dict[str, Any] | None:
        return ensure_safe_mapping(value, info.field_name)


class student_response(student_base):
    """Schema for student response."""

    model_config = ConfigDict(from_attributes=True)

    id: UUID
    created_at: datetime | None = None
    updated_at: datetime | None = None
