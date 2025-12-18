"""Clients package for external service integrations.

This package contains client modules for integrating with external services
and validating API contracts.

Modules:
    - **testing**: Contract validation client for Testing Service API
"""

from app.clients.testing import ContractManager, ContractValidationError, SchemaLoader

__all__ = ["ContractManager", "ContractValidationError", "SchemaLoader"]
