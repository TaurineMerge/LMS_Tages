"""Clients package for external service integrations.

This package contains client modules for integrating with external services
and validating API contracts.

Modules:
    - **testing**: Contract validation client for Testing Service API
"""

from app.clients.testing_client import TestingClient
from app.clients.validation.contract_manager import ContractManager, ContractValidationError
from app.clients.validation.schema_loader import SchemaLoader

__all__ = [
    "ContractManager",
    "ContractValidationError",
    "SchemaLoader",
    "TestingClient",
]
