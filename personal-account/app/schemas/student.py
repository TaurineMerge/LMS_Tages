"""Student schemas."""
from datetime import date, datetime
from uuid import UUID
from typing import Any

from pydantic import BaseModel, EmailStr, Field, ConfigDict


class StudentBase(BaseModel):
    """Base student schema with common fields."""
    
    name: str = Field(..., min_length=1, max_length=100, description="Имя студента")
    surname: str = Field(..., min_length=1, max_length=100, description="Фамилия студента")
    birth_date: date | None = Field(None, description="Дата рождения")
    avatar: str | None = Field(None, max_length=500, description="URL аватара")
    contacts: dict[str, Any] | None = Field(default_factory=dict, description="Контактная информация")
    email: EmailStr = Field(..., description="Email адрес")
    phone: str | None = Field(None, max_length=20, description="Номер телефона")


class StudentCreate(StudentBase):
    """Schema for creating a student."""
    pass


class StudentUpdate(BaseModel):
    """Schema for updating a student. All fields are optional."""
    
    model_config = ConfigDict(extra="forbid")
    
    name: str | None = Field(None, min_length=1, max_length=100)
    surname: str | None = Field(None, min_length=1, max_length=100)
    birth_date: date | None = None
    avatar: str | None = Field(None, max_length=500)
    contacts: dict[str, Any] | None = None
    email: EmailStr | None = None
    phone: str | None = Field(None, max_length=20)


class StudentResponse(StudentBase):
    """Schema for student response."""
    
    model_config = ConfigDict(from_attributes=True)
    
    id: UUID
    created_at: datetime | None = None
    updated_at: datetime | None = None
