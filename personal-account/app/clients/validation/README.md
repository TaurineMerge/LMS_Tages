# Testing Client - Contract Validation

High-performance API contract validation using fastjsonschema.

## 📋 Overview

This package provides comprehensive tools for validating API contracts against JSON Schema definitions. It's designed for testing API responses and requests, ensuring data integrity in microservice communication, and validating data pipelines.

### Key Features

- ⚡ **High Performance**: fastjsonschema compiled validators (10-100x faster than jsonschema)
- 🎯 **Schema Versioning**: Support for v1, v2, and automatic "latest" version
- 💾 **Smart Caching**: Compiled validators cached for optimal performance
- 📊 **Rich Error Reporting**: Structured error objects with JSON paths and details
- 🔄 **Async First**: Built for modern async/await workflows
- 🔒 **Type Safe**: Full type hints for excellent IDE support

## 🚀 Quick Start

### Installation

```bash
# Already included in personal-account dependencies
poetry install
```

### Basic Usage

```python
from app.clients.testing import ContractManager, ContractValidationError

# Initialize manager
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

## 📚 API Reference

### ContractManager

Main interface for contract validation.

```python
manager = ContractManager(
    schemas_dir=None,        # Optional: custom schemas directory
    cache_size=100,          # Max cached validators
    enable_caching=True      # Enable/disable caching
)
```

#### Methods

**`validate(data, contract_name, version="latest")`**
Generic validation method for any contract.

```python
validated = await manager.validate(
    data={"student_id": "...", ...},
    contract_name="user_stats",
    version="v1"
)
```

**`validate_user_stats(data)`**
Validate user statistics (convenience method).

**`validate_attempt_detail(data)`**
Validate test attempt details (convenience method).

**`validate_attempts_list(data)`**
Validate attempts collection (convenience method).

**`clear_cache()`**
Clear all cached validators (useful during development).

### SchemaLoader

Handles loading and caching of JSON schemas.

```python
loader = SchemaLoader(schemas_dir=None)
schema = await loader.load("user_stats", "v1")
```

### ContractValidationError

Exception raised on validation failures with detailed error information.

```python
try:
    await manager.validate_user_stats(invalid_data)
except ContractValidationError as e:
    print(e.message)          # Human-readable summary
    print(len(e.errors))      # Number of errors
    for err in e.errors:
        print(err['path'])    # JSON path to invalid field
        print(err['message']) # Error description
        print(err['rule'])    # Validation rule that failed
```

## 📂 Directory Structure

```
app/clients/testing/
├── __init__.py           # Package exports and documentation
├── contract_manager.py   # Main validation logic (fastjsonschema)
├── schema_loader.py      # Async schema loading and caching
├── examples.py           # Usage examples and patterns
├── README.md             # This file
└── schemas/              # JSON Schema definitions
    ├── user_stats/
    │   └── v1.json
    ├── attempt_detail/
    │   └── v1.json
    └── attempts_list/
        └── v1.json
```

## 📖 Examples

### Error Handling

```python
invalid_data = {
    "student_id": 12345,      # Should be UUID string
    "total_attempts": -1,     # Should be >= 0
}

try:
    await manager.validate_user_stats(invalid_data)
except ContractValidationError as e:
    for i, error in enumerate(e.errors, 1):
        print(f"Error {i}:")
        print(f"  Path: {error['path']}")
        print(f"  Message: {error['message']}")
        print(f"  Rule: {error['rule']}")
```

### Using Specific Schema Versions

```python
# Validate with explicit version
validated = await manager.validate(data, "user_stats", version="v1")

# Validate with "latest" (default - picks highest version)
validated = await manager.validate(data, "user_stats", version="latest")
```

### Custom Schemas Directory

```python
from pathlib import Path

manager = ContractManager(
    schemas_dir=Path("/custom/schemas"),
    cache_size=200
)
```

### Disabling Cache (Development)

```python
# Disable caching when schemas change frequently
manager = ContractManager(enable_caching=False)

# Or clear cache manually
manager.clear_cache()
```

## 🔧 Schema Format

Schemas are JSON Schema Draft-07 compliant:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "user_stats_v1",
  "title": "UserStats",
  "description": "User statistics contract",
  "type": "object",
  "required": ["student_id", "total_attempts"],
  "properties": {
    "student_id": {
      "type": "string",
      "format": "uuid",
      "description": "Unique student identifier"
    },
    "total_attempts": {
      "type": "integer",
      "minimum": 0,
      "description": "Total number of test attempts"
    }
  }
}
```

## 🧪 Testing

```python
import pytest
from app.clients.testing import ContractManager, ContractValidationError

@pytest.fixture
def manager():
    return ContractManager()

@pytest.mark.asyncio
async def test_valid_user_stats(manager):
    data = {
        "student_id": "550e8400-e29b-41d4-a716-446655440000",
        "total_attempts": 5,
        "passed_attempts": 3,
        "completion_rate": 60.0,
        "average_score": 75.5,
        "last_attempt_at": "2025-12-17T14:30:00Z"
    }
    result = await manager.validate_user_stats(data)
    assert result == data

@pytest.mark.asyncio
async def test_invalid_user_stats(manager):
    invalid_data = {"student_id": 12345}  # Wrong type

    with pytest.raises(ContractValidationError) as exc_info:
        await manager.validate_user_stats(invalid_data)

    assert len(exc_info.value.errors) > 0
```

## 📈 Performance

fastjsonschema compiles JSON Schema into optimized Python code:

| Library        | Validations/sec |
|---------------|-----------------|
| fastjsonschema | ~50,000        |
| jsonschema     | ~500           |

Benchmark: Simple object validation, cached validator.

## 🔗 Dependencies

- **fastjsonschema**: Fast JSON Schema validation (compiled validators)
- **aiofiles**: Async file operations for schema loading

## 📝 License

MIT License - Part of LMS_Tages project.
