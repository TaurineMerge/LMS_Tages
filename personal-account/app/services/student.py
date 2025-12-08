"""Student service layer."""
from uuid import UUID

from app.repositories.student import student_repository
from app.schemas.student import student_create, student_update, student_response
from app.schemas.common import paginated_response
from app.exceptions import not_found_error, conflict_error
from app.telemetry import traced


class student_service:
    """Service for student business logic."""
    
    def __init__(self):
        self.repository = student_repository
    
    @traced()
    async def get_students(
        self, 
        page: int = 1, 
        limit: int = 20
    ) -> paginated_response[student_response]:
        """Get paginated list of students."""
        students, total = await self.repository.get_paginated(page, limit)
        
        return paginated_response(
            data=[student_response(**s) for s in students],
            total=total,
            page=page,
            limit=limit
        )
    
    @traced()
    async def get_student(self, student_id: UUID) -> student_response:
        """Get student by ID."""
        student = await self.repository.get_by_id(student_id)
        
        if not student:
            raise not_found_error("Student", str(student_id))
        
        return student_response(**student)
    
    @traced()
    async def create_student(self, data: student_create) -> student_response:
        """Create a new student."""
        # Check for duplicate email
        if await self.repository.email_exists(data.email):
            raise conflict_error(f"Student with email '{data.email}' already exists")
        
        student = await self.repository.create(data.model_dump())
        
        if not student:
            raise Exception("Failed to create student")
        
        return student_response(**student)
    
    @traced()
    async def update_student(
        self, 
        student_id: UUID, 
        data: student_update
    ) -> student_response:
        """Update student by ID."""
        # Check if student exists
        existing = await self.repository.get_by_id(student_id)
        if not existing:
            raise not_found_error("Student", str(student_id))
        
        # Check for email conflict if email is being updated
        if data.email and await self.repository.email_exists(data.email, student_id):
            raise conflict_error(f"Email '{data.email}' is already in use")
        
        student = await self.repository.update(
            student_id, 
            data.model_dump(exclude_unset=True)
        )
        
        if not student:
            raise Exception("Failed to update student")
        
        return student_response(**student)
    
    @traced()
    async def delete_student(self, student_id: UUID) -> bool:
        """Delete student by ID."""
        if not await self.repository.exists(student_id):
            raise not_found_error("Student", str(student_id))
        
        return await self.repository.delete(student_id)


# Singleton instance
student_service = student_service()
