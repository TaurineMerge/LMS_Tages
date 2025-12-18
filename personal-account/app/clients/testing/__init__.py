"""Testing client package for API contract validation.

This package provides comprehensive tools for validating API contracts using JSON Schema
and the fastjsonschema library. It supports versioned schemas, fast validation with
caching, and detailed error reporting for API testing.

Core Components:
    - **ContractManager**: High-level API for contract validation
    - **SchemaLoader**: Async schema loading with caching
    - **ContractValidationError**: Rich exception with detailed error info

Quick Start:

    ```python
    from app.clients.testing import ContractManager, ContractValidationError

    manager = ContractManager()

    # Validate user statistics
    try:
        data = {
            "student_id": "550e8400-e29b-41d4-a716-446655440000",
            "total_attempts": 5,
            "passed_attempts": 3,
            "completion_rate": 60.0,
            "average_score": 75.5,
            "last_attempt_at": "2025-12-17T14:30:00Z"
        }
        validated = await manager.validate_user_stats(data)
        print("✓ Valid!")
    except ContractValidationError as e:
        print(f"✗ {e.message}")
        for error in e.errors:
            print(f"  [{error['path']}] {error['message']}")
    ```

Features:
    - **High Performance**: fastjsonschema compiled validators
    - **Schema Versioning**: Support for v1, v2, and automatic "latest" version
    - **Smart Caching**: Compiled validators cached for repeated validations
    - **Rich Errors**: Structured error objects with JSON paths and details
    - **Async First**: Built for modern async/await workflows
    - **Type Safe**: Full type hints for excellent IDE support

Supported Contracts:
    - `user_stats`: User progress and statistics
    - `attempt_detail`: Individual test attempt details
    - `attempts_list`: Collections of test attempts with pagination

Package Structure:
    ```
    app/clients/testing/
    ├── __init__.py           # Package exports (this file)
    ├── contract_manager.py   # Main validation logic with fastjsonschema
    ├── schema_loader.py      # Async schema loading and caching
    ├── examples.py           # Usage examples and patterns
    └── schemas/              # JSON Schema definitions
        ├── user_stats/
        │   └── v1.json
        ├── attempt_detail/
        │   └── v1.json
        └── attempts_list/
            └── v1.json
    ```

See Also:
    - JSON Schema documentation: https://json-schema.org/
    - fastjsonschema: https://pypi.org/project/fastjsonschema/

Examples:
    See `examples.py` for detailed usage patterns and integration examples.
"""

from .contract_manager import ContractManager, ContractValidationError
from .schema_loader import SchemaLoader

__all__ = [
    "ContractManager",
    "ContractValidationError",
    "SchemaLoader",
]

__version__ = "1.0.0"
