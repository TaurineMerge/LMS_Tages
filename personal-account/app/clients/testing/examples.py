"""Usage examples and integration patterns for the testing client.

This module provides comprehensive examples demonstrating how to use the ContractManager
for API contract validation in various scenarios including API testing, data pipelines,
and integration tests.

Examples cover:
    - Basic validation workflows
    - Error handling patterns
    - Integration with FastAPI endpoints
    - Batch validation
    - Custom validation scenarios

Usage:

    ```python
    # Import the examples module
    from app.clients.testing import examples

    # Run example validations
    await examples.basic_validation_example()
    await examples.error_handling_example()
    await examples.batch_validation_example()
    ```

See Also:
    - `ContractManager`: Main validation interface
    - `ContractValidationError`: Exception with detailed errors
    - `SchemaLoader`: Schema loading utilities
"""

import logging
from typing import Any

from .contract_manager import ContractManager, ContractValidationError

logger = logging.getLogger(__name__)


async def basic_validation_example() -> None:
    """Demonstrate basic contract validation workflow.

    Shows how to:
    - Initialize ContractManager
    - Validate data against user_stats contract
    - Handle successful validation

    Example output:
        ```
        ✓ User statistics validated successfully
        Validated data: {'student_id': '...', 'total_attempts': 5, ...}
        ```
    """
    manager = ContractManager()

    # Sample user statistics data
    user_stats = {
        "student_id": "550e8400-e29b-41d4-a716-446655440000",
        "total_attempts": 5,
        "passed_attempts": 3,
        "completion_rate": 60.0,
        "average_score": 75.5,
        "last_attempt_at": "2025-12-17T14:30:00Z",
    }

    try:
        validated = await manager.validate_user_stats(user_stats)
        logger.info("✓ User statistics validated successfully")
        logger.info(f"Validated data: {validated}")
    except ContractValidationError as e:
        logger.error(f"✗ Validation failed: {e.message}")


async def error_handling_example() -> None:
    """Demonstrate error handling with detailed error information.

    Shows how to:
    - Handle validation failures gracefully
    - Extract detailed error information
    - Format errors for logging or user display

    Example output:
        ```
        ✗ Validation failed for user_stats
        Total errors: 2
        Error 1:
          Path: student_id
          Issue: must be string
          Keyword: type
        Error 2:
          Path: total_attempts
          Issue: must be >= 0
          Keyword: minimum
        ```
    """
    manager = ContractManager()

    # Invalid data (wrong types, missing fields)
    invalid_stats = {
        "student_id": 12345,  # Should be UUID string
        "total_attempts": -1,  # Should be >= 0
        # Missing required fields
    }

    try:
        await manager.validate_user_stats(invalid_stats)
    except ContractValidationError as e:
        logger.error("✗ Validation failed for user_stats")
        logger.error(f"Total errors: {len(e.errors)}")

        for i, error in enumerate(e.errors, 1):
            logger.error(f"Error {i}:")
            logger.error(f"  Path: {error['path']}")
            logger.error(f"  Issue: {error['message']}")
            logger.error(f"  Rule: {error['rule']}")


async def batch_validation_example() -> None:
    """Demonstrate batch validation of multiple records.

    Shows how to:
    - Validate multiple records efficiently
    - Collect validation results
    - Handle partial failures in batch processing

    Example output:
        ```
        Processing batch of 3 records...
        ✓ Record 1: Valid
        ✗ Record 2: Invalid (2 errors)
        ✓ Record 3: Valid
        Batch results: 2/3 valid
        ```
    """
    manager = ContractManager()

    # Batch of user statistics records
    batch_data = [
        {
            "student_id": "550e8400-e29b-41d4-a716-446655440000",
            "total_attempts": 10,
            "passed_attempts": 8,
            "completion_rate": 80.0,
            "average_score": 85.5,
            "last_attempt_at": "2025-12-17T10:00:00Z",
        },
        {
            "student_id": "invalid-uuid",  # Invalid
            "total_attempts": "not-a-number",  # Invalid
            "passed_attempts": 5,
            "completion_rate": 100.0,
            "average_score": 90.0,
            "last_attempt_at": "2025-12-17T11:00:00Z",
        },
        {
            "student_id": "660e8400-e29b-41d4-a716-446655440001",
            "total_attempts": 3,
            "passed_attempts": 3,
            "completion_rate": 100.0,
            "average_score": 95.0,
            "last_attempt_at": "2025-12-17T12:00:00Z",
        },
    ]

    logger.info(f"Processing batch of {len(batch_data)} records...")

    results: list[tuple[int, bool, str]] = []

    for i, record in enumerate(batch_data, 1):
        try:
            await manager.validate_user_stats(record)
            logger.info(f"✓ Record {i}: Valid")
            results.append((i, True, "Valid"))
        except ContractValidationError as e:
            error_count = len(e.errors)
            logger.warning(f"✗ Record {i}: Invalid ({error_count} errors)")
            results.append((i, False, e.message))

    valid_count = sum(1 for _, is_valid, _ in results if is_valid)
    logger.info(f"Batch results: {valid_count}/{len(batch_data)} valid")


async def fastapi_integration_example(data: dict[str, Any]) -> dict[str, Any]:
    """Demonstrate integration with FastAPI endpoint.

    Shows how to:
    - Validate request/response data in FastAPI
    - Return validation errors as HTTP 422
    - Use ContractManager in endpoint handlers

    Args:
        data: Request data to validate.

    Returns:
        Validated data.

    Raises:
        HTTPException: 422 if validation fails.

    Example FastAPI usage:

        ```python
        from fastapi import FastAPI, HTTPException
        from app.clients.testing.examples import fastapi_integration_example

        app = FastAPI()

        @app.post("/api/v1/statistics/user")
        async def create_user_stats(data: dict):
            try:
                validated = await fastapi_integration_example(data)
                # Process validated data...
                return {"status": "success", "data": validated}
            except HTTPException:
                raise
        ```
    """
    from fastapi import HTTPException

    manager = ContractManager()

    try:
        validated = await manager.validate_user_stats(data)
        logger.info("✓ Request validated successfully")
        return validated

    except ContractValidationError as e:
        logger.warning(f"✗ Request validation failed: {e.message}")

        # Format errors for HTTP response
        error_details = [{"field": err["path"], "message": err["message"], "type": err["keyword"]} for err in e.errors]

        raise HTTPException(
            status_code=422,
            detail={"message": "Validation failed", "errors": error_details},
        ) from e


async def custom_contract_example() -> None:
    """Demonstrate validation with custom contract and version.

    Shows how to:
    - Use generic validate() method
    - Specify contract name and version explicitly
    - Work with different schema versions

    Example output:
        ```
        ✓ Custom contract validated: attempt_detail:v1
        ```
    """
    manager = ContractManager()

    attempt_data = {
        "attempt_id": "770e8400-e29b-41d4-a716-446655440002",
        "student_id": "550e8400-e29b-41d4-a716-446655440000",
        "test_id": "880e8400-e29b-41d4-a716-446655440003",
        "score": 85.5,
        "completed": True,
        "started_at": "2025-12-17T14:30:00Z",
        "completed_at": "2025-12-17T15:00:00Z",
        "answers": [],
    }

    try:
        # Use generic validate() with explicit contract and version
        validated = await manager.validate(
            data=attempt_data,
            contract_name="attempt_detail",
            version="v1",
        )
        logger.info("✓ Custom contract validated: attempt_detail:v1")
        logger.debug(f"Validated: {validated}")
    except ContractValidationError as e:
        logger.error(f"✗ Validation failed: {e.message}")


async def caching_example() -> None:
    """Demonstrate validator caching for performance.

    Shows how to:
    - Benefit from automatic validator caching
    - Clear cache when schemas change
    - Configure caching behavior

    Example output:
        ```
        First validation: 12.5ms (compile + validate)
        Second validation: 0.8ms (cached)
        After cache clear: 11.2ms (recompile)
        ```
    """
    import time

    # Initialize with caching enabled (default)
    manager = ContractManager(enable_caching=True, cache_size=50)

    sample_data = {
        "student_id": "550e8400-e29b-41d4-a716-446655440000",
        "total_attempts": 5,
        "passed_attempts": 3,
        "completion_rate": 60.0,
        "average_score": 75.5,
        "last_attempt_at": "2025-12-17T14:30:00Z",
    }

    # First validation (compiles schema)
    start = time.perf_counter()
    await manager.validate_user_stats(sample_data)
    first_time = (time.perf_counter() - start) * 1000
    logger.info(f"First validation: {first_time:.1f}ms (compile + validate)")

    # Second validation (uses cached validator)
    start = time.perf_counter()
    await manager.validate_user_stats(sample_data)
    second_time = (time.perf_counter() - start) * 1000
    logger.info(f"Second validation: {second_time:.1f}ms (cached)")

    # Clear cache (useful when schemas change during development)
    manager.clear_cache()

    # Third validation (recompiles)
    start = time.perf_counter()
    await manager.validate_user_stats(sample_data)
    third_time = (time.perf_counter() - start) * 1000
    logger.info(f"After cache clear: {third_time:.1f}ms (recompile)")


async def run_all_examples() -> None:
    """Run all example functions for demonstration.

    Executes all examples in sequence to show the full range of
    ContractManager capabilities.

    Example:

        ```python
        from app.clients.testing.examples import run_all_examples

        # Run all examples
        await run_all_examples()
        ```
    """
    logger.info("=" * 60)
    logger.info("Running Contract Validation Examples")
    logger.info("=" * 60)

    examples = [
        ("Basic Validation", basic_validation_example),
        ("Error Handling", error_handling_example),
        ("Batch Validation", batch_validation_example),
        ("Custom Contract", custom_contract_example),
        ("Caching Performance", caching_example),
    ]

    for name, example_func in examples:
        logger.info(f"\n--- {name} ---")
        try:
            await example_func()
        except (ContractValidationError, ValueError, OSError) as e:
            logger.error(f"Example failed: {e}")

    logger.info("\n" + "=" * 60)
    logger.info("All examples completed")
    logger.info("=" * 60)


# Example data templates for testing
EXAMPLE_USER_STATS = {
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "total_attempts": 5,
    "passed_attempts": 3,
    "completion_rate": 60.0,
    "average_score": 75.5,
    "last_attempt_at": "2025-12-17T14:30:00Z",
}

EXAMPLE_ATTEMPT_DETAIL = {
    "attempt_id": "770e8400-e29b-41d4-a716-446655440002",
    "student_id": "550e8400-e29b-41d4-a716-446655440000",
    "test_id": "880e8400-e29b-41d4-a716-446655440003",
    "score": 85.5,
    "completed": True,
    "started_at": "2025-12-17T14:30:00Z",
    "completed_at": "2025-12-17T15:00:00Z",
    "answers": [],
}

EXAMPLE_ATTEMPTS_LIST = {
    "attempts": [EXAMPLE_ATTEMPT_DETAIL],
    "total": 1,
    "page": 1,
    "page_size": 20,
}
