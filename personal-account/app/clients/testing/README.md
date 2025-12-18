# Testing Client - Contract Validation

High-performance API contract validation using AJV (Another JSON Schema Validator) via pyajv.

## 📋 Overview

This package provides comprehensive tools for validating API contracts against JSON Schema definitions. It's designed for testing API responses and requests, ensuring data integrity in microservice communication, and validating data pipelines.

### Key Features

- ⚡ **High Performance**: AJV-powered validation with compiled schemas
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
        print(err['keyword']) # AJV keyword that failed
```

## 📂 Directory Structure

```
app/clients/testing/
├── __init__.py           # Package exports and documentation
├── contract_manager.py   # Main validation logic (AJV-based)
├── schema_loader.py      # Async schema loading and caching
├── examples.py           # Usage examples and patterns
├── README.md            # This file
└── schemas/             # JSON Schema definitions
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
        print(f"  Issue: {error['message']}")
        print(f"  Keyword: {error['keyword']}")
```

### Batch Validation

```python
batch_data = [record1, record2, record3]
results = []

for i, record in enumerate(batch_data, 1):
    try:
        await manager.validate_user_stats(record)
        results.append((i, True, "Valid"))
    except ContractValidationError as e:
        results.append((i, False, e.message))

valid_count = sum(1 for _, is_valid, _ in results if is_valid)
print(f"Results: {valid_count}/{len(batch_data)} valid")
```

### FastAPI Integration

```python
from fastapi import FastAPI, HTTPException

app = FastAPI()
manager = ContractManager()

@app.post("/api/v1/statistics/user")
async def create_user_stats(data: dict):
    try:
        validated = await manager.validate_user_stats(data)
        # Process validated data...
        return {"status": "success", "data": validated}
    except ContractValidationError as e:
        raise HTTPException(
            status_code=422,
            detail={
                "message": "Validation failed",
                "errors": [
                    {
                        "field": err["path"],
                        "message": err["message"],
                        "type": err["keyword"]
                    }
                    for err in e.errors
                ]
            }
        )
```

### Performance & Caching

```python
# Caching enabled by default
manager = ContractManager(enable_caching=True, cache_size=50)

# First validation: compiles schema
await manager.validate_user_stats(data)  # ~12ms

# Second validation: uses cached validator
await manager.validate_user_stats(data)  # ~0.8ms

# Clear cache when schemas change
manager.clear_cache()
```

## 🛠 Development

### Adding New Schemas

1. Create directory: `schemas/new_contract/`
2. Add schema: `schemas/new_contract/v1.json`
3. Add convenience method to `ContractManager`:

```python
async def validate_new_contract(self, data: dict[str, Any]) -> dict[str, Any]:
    """Validate NewContract (v1)."""
    return await self.validate(data, "new_contract", "v1")
```

### Running Examples

```python
from app.clients.testing.examples import run_all_examples

# Run all example scenarios
await run_all_examples()
```

### Testing

```python
import pytest
from app.clients.testing import ContractManager, ContractValidationError

@pytest.mark.asyncio
async def test_valid_user_stats():
    manager = ContractManager()
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
async def test_invalid_user_stats():
    manager = ContractManager()
    with pytest.raises(ContractValidationError) as exc_info:
        await manager.validate_user_stats({"invalid": "data"})
    assert len(exc_info.value.errors) > 0
```

## 🔍 Supported Contracts

### user_stats (v1)

User progress and statistics.

**Fields:**
- `student_id` (uuid): Student identifier
- `total_attempts` (int): Total attempts count
- `passed_attempts` (int): Passed attempts count
- `completion_rate` (float): Completion percentage (0-100)
- `average_score` (float): Average score
- `last_attempt_at` (datetime): ISO 8601 timestamp

### attempt_detail (v1)

Detailed information about a test attempt.

**Fields:**
- `attempt_id` (uuid): Attempt identifier
- `student_id` (uuid): Student identifier
- `test_id` (uuid): Test identifier
- `score` (float): Final score
- `completed` (bool): Completion status
- `started_at` (datetime): Start timestamp
- `completed_at` (datetime?): Completion timestamp (nullable)
- `answers` (array): List of answer objects

### attempts_list (v1)

Collection of attempts with pagination.

**Fields:**
- `attempts` (array): List of attempt objects
- `total` (int): Total count
- `page` (int): Current page number
- `page_size` (int): Items per page

## 📝 Notes

- **Performance**: AJV validators are compiled once and cached. First validation is slower (~10-15ms), subsequent validations are very fast (~0.5-1ms).
- **Schema Changes**: Call `manager.clear_cache()` after updating schemas during development.
- **Error Handling**: Always catch `ContractValidationError` and handle errors appropriately in production.
- **Version Management**: Use explicit versions in production (`v1`, `v2`) rather than `"latest"`.

## 🔗 See Also

- [JSON Schema Documentation](https://json-schema.org/)
- [AJV Documentation](https://ajv.js.org/)
- [pyajv Python Package](https://pypi.org/project/pyajv/)
- [FastAPI Documentation](https://fastapi.tiangolo.com/)

## 📄 License

Part of the LMS_Tages personal-account service.
