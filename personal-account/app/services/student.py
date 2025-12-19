"""Student service layer."""

from uuid import UUID

# FastAPI
from fastapi import HTTPException
from fastapi.concurrency import run_in_threadpool

from app.database import execute, fetch_one

# App modules
from app.exceptions import conflict_error, not_found_error
from app.repositories.student import student_repository
from app.schemas.common import paginated_response
from app.schemas.student import student_create, student_response, student_update
from app.services.keycloak import keycloak_service
from app.telemetry import traced


class student_service:
    """Service for student business logic."""

    def __init__(self):
        self.repository = student_repository

    @traced()
    async def get_students(self, page: int = 1, limit: int = 20) -> paginated_response[student_response]:
        """Get paginated list of students."""
        students, total = await self.repository.get_paginated(page, limit)

        return paginated_response(data=[student_response(**s) for s in students], total=total, page=page, limit=limit)

    @traced()
    async def get_student(self, student_id: UUID) -> student_response:
        """Get student by ID."""
        student = await self.repository.get_by_id(student_id)

        if not student:
            raise not_found_error("Student", str(student_id))

        return student_response(**student)

    # @traced()
    # async def get_student_by_email(self, email: str) -> student_response:
    #     """Get student by email."""
    #     student = await self.repository.get_by_email(email)

    #     if not student:
    #         raise not_found_error("Student", email)

    #     return student_response(**student)

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
    async def update_student(self, student_id: UUID, data: student_update):
        """Update student in both Keycloak and DB atomically."""
        # Получаем старые данные Keycloak, чтобы откатить при ошибке в БД
        old_keycloak_data = await run_in_threadpool(keycloak_service.get_user_data, str(student_id))

        # Открываем транзакцию через вашу систему
        from app.database import _engine

        if _engine is None:
            raise RuntimeError("Database engine is not initialized")

        async with _engine.begin() as conn:  # ← ваша транзакция
            try:
                # 1. Обновляем Keycloak
                keycloak_payload = {}
                if data.name is not None:
                    keycloak_payload["firstName"] = data.name
                if data.surname is not None:
                    keycloak_payload["lastName"] = data.surname
                if data.email is not None:
                    keycloak_payload["email"] = data.email
                if data.username is not None:
                    keycloak_payload["username"] = data.username

                await run_in_threadpool(keycloak_service.update_user_data, str(student_id), keycloak_payload)

                # 2. Проверяем конфликты в БД (например, email)
                if data.email:
                    # Замените на вашу проверку (например, через репозиторий)
                    existing = await self.repository.email_exists(data.email, student_id, conn=conn)
                    if existing:
                        raise conflict_error(f"Email '{data.email}' is already in use")

                # 3. Обновляем БД
                student = await self.repository.update(student_id, data.model_dump(exclude_unset=True), conn=conn)

                if not student:
                    raise Exception("Failed to update student in DB")

                # conn.commit() — автоматически вызывается при выходе из `async with _engine.begin()`

            except Exception as e:
                # conn.rollback() — автоматически вызывается при выходе из `async with` при ошибке

                # Пытаемся откатить изменения в Keycloak
                try:
                    await run_in_threadpool(keycloak_service.update_user_data, str(student_id), old_keycloak_data)
                except Exception as ke:
                    print(f"Failed to rollback Keycloak changes: {ke}")
                    raise HTTPException(status_code=500, detail="Rollback failed for Keycloak.")

                raise e  # пробрасываем ошибку дальше

    @traced()
    async def delete_student(self, student_id: UUID) -> bool:
        """Delete student by ID."""
        if not await self.repository.exists(student_id):
            raise not_found_error("Student", str(student_id))

        return await self.repository.delete(student_id)


# Singleton instance
student_service = student_service()
