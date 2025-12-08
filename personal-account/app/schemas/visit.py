"""Visit schemas."""
from uuid import UUID

from pydantic import BaseModel, ConfigDict


class visit_base(BaseModel):
    """Base visit schema."""
    
    student_id: UUID
    lesson_id: UUID


class visit_create(visit_base):
    """Schema for creating a visit."""
    pass


class visit_response(visit_base):
    """Schema for visit response."""
    
    model_config = ConfigDict(from_attributes=True)
    
    id: UUID
