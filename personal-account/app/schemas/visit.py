"""Visit schemas."""
from uuid import UUID

from pydantic import BaseModel, ConfigDict


class VisitBase(BaseModel):
    """Base visit schema."""
    
    student_id: UUID
    lesson_id: UUID


class VisitCreate(VisitBase):
    """Schema for creating a visit."""
    pass


class VisitResponse(VisitBase):
    """Schema for visit response."""
    
    model_config = ConfigDict(from_attributes=True)
    
    id: UUID
