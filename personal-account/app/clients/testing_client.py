"""HTTP Testing client that talks to the external testing service.

This client implements three read-only endpoints used by our ingestion pipeline:

- GET /internal/users/{user_id}/stats -> returns aggregates (user stats)
- GET /internal/users/{user_id}/attempts -> returns list of attempts (light items)
- GET /internal/attempts/{attempt_id} -> returns detailed attempt info

The client validates responses using :class:`ContractManager` and persists
the raw payloads into our integration schema using the
:mod:`app.repositories.integration` repository.

Instrumentation
---------------
All public methods are decorated with ``@traced()`` so calls, parameters,
returned values and exceptions are visible in Jaeger (when OpenTelemetry is
enabled). The implementation favors defensive imports so documentation
generation and light-weight introspection do not require full runtime
dependencies.

Example
-------
```py
from app.clients.testing import TestingClient

client = TestingClient(base_url="https://testing.example")
data = await client.fetch_user_stats("550e8400-e29b-41d4-a716-446655440000")
```
"""

from __future__ import annotations

import logging
from typing import Any
from uuid import UUID

import httpx
import settings
from contract_manager import ContractManager, ContractValidationError

from app.repositories.integration import integration_repository
from app.telemetry import traced

logger = logging.getLogger(__name__)


class TestingClient:
    """Async client for the external testing service.

    This client is intentionally small and focused: it hits three endpoints,
    validates responses with :class:`ContractManager`, and stores raw payloads
    via the :mod:`app.repositories.integration` repository for downstream
    processing.

    Parameters
    ----------
    base_url: str
        Base URL of the testing service (no trailing slash preferred).
    timeout_seconds: int
        Per-request timeout passed to httpx.
    contract_manager: ContractManager | None
        Optional ContractManager instance. If omitted, a default will be
        created which loads schemas from the package.
    repo: integration_repository | None
        Repository instance used to persist raw payloads. Defaults to the
        project singleton ``integration_repository``.
    """

    def __init__(
        self,
        base_url: str,
        timeout_seconds: int = 10,
        contract_manager: ContractManager | None = None,
        repo: Any | None = None,
    ) -> None:
        self.base_url = settings.TESTING_BASE_URL
        self.timeout = timeout_seconds
        self.contract_manager = contract_manager or ContractManager()
        self.repo = repo or integration_repository

    @traced()
    async def fetch_user_stats(self, user_id: UUID) -> dict[str, Any]:
        """Fetch user aggregates and persist the raw payload.

        - Calls GET /internal/users/{user_id}/stats
        - Validates response with ContractManager.validate_user_stats
        - Persists raw payload to integration.raw_user_stats

        Returns the validated payload (as dict).
        """
        url = f"{self.base_url}/internal/users/{user_id}/stats"
        logger.debug("Fetching user stats: %s", url)

        # import httpx locally so documentation generation and light imports
        # don't require httpx to be installed at import-time

        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.get(url)
            resp.raise_for_status()
            payload = resp.json()

        # Validate contract
        try:
            validated = await self.contract_manager.validate_user_stats(payload)
        except ContractValidationError as e:
            logger.warning("User stats contract validation failed for %s: %s", user_id, e)
            # Still persist raw payload for debugging/inspection
            await self.repo.insert_raw_user_stats(UUID(str(user_id)), payload)
            raise

        # Persist raw payload for processing
        await self.repo.insert_raw_user_stats(UUID(str(user_id)), payload)
        logger.debug("User stats fetched and persisted for %s", user_id)
        return validated  # type: ignore[return-value]

    @traced()
    async def fetch_user_attempts(self, user_id: UUID) -> list[dict[str, Any]]:
        """Fetch lightweight attempts list and persist each attempt as raw_attempt.

        - Calls GET /internal/users/{user_id}/attempts
        - Validates the list using ContractManager.validate_attempts_list
        - Persists each attempt into integration.raw_attempts using
          external_attempt_id = attempt['attempt_id']

        Returns the validated list of attempts.
        """
        url = f"{self.base_url}/internal/users/{user_id}/attempts"
        logger.debug("Fetching attempts list: %s", url)

        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.get(url)
            resp.raise_for_status()
            payload = resp.json()

            # Validate contract
            try:
                validated = await self.contract_manager.validate_attempts_list(payload)
            except ContractValidationError as e:
                logger.warning("Attempts list validation failed for %s: %s", user_id, e)
                # Persist the entire payload as user stats fallback
                await self.repo.insert_raw_user_stats(UUID(str(user_id)), payload)
                raise

            # Persist each attempt as raw_attempt
            for item in validated:
                external_attempt_id = UUID(str(item.get("attempt_id")))
                test_id = None
                if item.get("test_id"):
                    test_id = UUID(str(item.get("test_id")))
                await self.repo.insert_raw_attempt(external_attempt_id, UUID(str(user_id)), test_id, item)

            logger.debug("Attempts list fetched and %d attempts persisted for %s", len(validated), user_id)
            return validated  # type: ignore[return-value]

    @traced()
    async def fetch_attempt_detail(self, attempt_id: UUID) -> dict[str, Any]:
        """Fetch attempt detail and persist it as a raw_attempt.

        - Calls GET /internal/attempts/{attempt_id}
        - Validates response with ContractManager.validate_attempt_detail
        - Persists raw attempt using external_attempt_id = attempt_id

        Returns the validated attempt detail.
        """
        url = f"{self.base_url}/internal/attempts/{attempt_id}"
        logger.debug("Fetching attempt detail: %s", url)

        async with httpx.AsyncClient(timeout=self.timeout) as client:
            resp = await client.get(url)
            resp.raise_for_status()
            payload = resp.json()

            # Validate contract
            try:
                validated = await self.contract_manager.validate_attempt_detail(payload)
            except ContractValidationError as e:
                logger.warning("Attempt detail validation failed for %s: %s", attempt_id, e)
                # Persist raw payload for inspection
                await self.repo.insert_raw_attempt(
                    UUID(str(attempt_id)), UUID(str(payload.get("student_id"))), None, payload
                )
                raise

            # Persist raw attempt
            student_id = UUID(str(validated.get("student_id"))) if validated.get("student_id") else None
            await self.repo.insert_raw_attempt(
                UUID(str(attempt_id)),
                student_id,
                UUID(str(validated.get("test_id"))) if validated.get("test_id") else None,
                validated,
            )

            logger.debug("Attempt detail fetched and persisted for %s", attempt_id)
            return validated  # type: ignore[return-value]


__all__ = ["TestingClient"]
