"""Student service layer."""
from uuid import UUID

from app.repositories.student import student_repository
from app.schemas.student import StudentCreate, StudentUpdate, StudentResponse
from app.schemas.common import PaginatedResponse
from app.exceptions import NotFoundError, ConflictError


class StudentService:
    """Service for student business logic."""
    
    def __init__(self):
        self.repository = student_repository
    
    async def get_students(
        self, 
        page: int = 1, 
        limit: int = 20
    ) -> PaginatedResponse[StudentResponse]:
        """Get paginated list of students."""
        students, total = await self.repository.get_paginated(page, limit)
        
        return PaginatedResponse(
            data=[StudentResponse(**s) for s in students],
            total=total,
            page=page,
            limit=limit
        )
    
    async def get_student(self, student_id: UUID) -> StudentResponse:
        """Get student by ID."""
        student = await self.repository.get_by_id(student_id)
        
        if not student:
            raise NotFoundError("Student", str(student_id))
        
        return StudentResponse(**student)
    
    async def create_student(self, data: StudentCreate) -> StudentResponse:
        """Create a new student."""
        # Check for duplicate email
        if await self.repository.email_exists(data.email):
            raise ConflictError(f"Student with email '{data.email}' already exists")
        
        student = await self.repository.create(data.model_dump())
        
        if not student:
            raise Exception("Failed to create student")
        
        return StudentResponse(**student)
    
    async def update_student(
        self, 
        student_id: UUID, 
        data: StudentUpdate
    ) -> StudentResponse:
        """Update student by ID."""
        # Check if student exists
        existing = await self.repository.get_by_id(student_id)
        if not existing:
            raise NotFoundError("Student", str(student_id))
        
        # Check for email conflict if email is being updated
        if data.email and await self.repository.email_exists(data.email, student_id):
            raise ConflictError(f"Email '{data.email}' is already in use")
        
        student = await self.repository.update(
            student_id, 
            data.model_dump(exclude_unset=True)
        )
        
        if not student:
            raise Exception("Failed to update student")
        
        return StudentResponse(**student)
    
    async def delete_student(self, student_id: UUID) -> bool:
        """Delete student by ID."""
        if not await self.repository.exists(student_id):
            raise NotFoundError("Student", str(student_id))
        
        return await self.repository.delete(student_id)


# Singleton instance
student_service = StudentService()
