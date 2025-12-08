"""Visit service layer."""
from uuid import UUID

from app.repositories.visit import visit_repository
from app.repositories.student import student_repository
from app.schemas.visit import visit_create, visit_response
from app.exceptions import not_found_error, conflict_error
from app.telemetry import traced


class visit_service:
    """Service for visit business logic."""
    
    def __init__(self):
        self.repository = visit_repository
        self.student_repository = student_repository
    
    @traced()
    async def get_visits(
        self,
        student_id: UUID | None = None,
        lesson_id: UUID | None = None
    ) -> list[visit_response]:
        """Get visits with optional filters."""
        visits = await self.repository.get_filtered(student_id, lesson_id)
        return [visit_response(**v) for v in visits]
    
    @traced()
    async def get_visit(self, visit_id: UUID) -> visit_response:
        """Get visit by ID."""
        visit = await self.repository.get_by_id(visit_id)
        
        if not visit:
            raise not_found_error("Visit", str(visit_id))
        
        return visit_response(**visit)
    
    @traced()
    async def create_visit(self, data: visit_create) -> visit_response:
        """Create a new visit record."""
        # Verify student exists
        if not await self.student_repository.exists(data.student_id):
            raise not_found_error("Student", str(data.student_id))
        
        # Note: Lesson validation would require cross-schema query
        # In production, you'd validate lesson_id exists in knowledge_base.lesson_d
        
        # Check for duplicate visit
        if await self.repository.visit_exists(data.student_id, data.lesson_id):
            raise conflict_error("Visit record already exists for this student and lesson")
        
        visit = await self.repository.create(data.model_dump())
        
        if not visit:
            # This shouldn't happen due to ON CONFLICT DO NOTHING
            raise conflict_error("Visit record already exists")
        
        return visit_response(**visit)
    
    @traced()
    async def delete_visit(self, visit_id: UUID) -> bool:
        """Delete visit by ID."""
        if not await self.repository.exists(visit_id):
            raise not_found_error("Visit", str(visit_id))
        
        return await self.repository.delete(visit_id)


# Singleton instance
visit_service = visit_service()
