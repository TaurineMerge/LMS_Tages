"""Contract validation manager for API contract testing.

This module provides high-performance contract validation using fastjsonschema library.
It offers fast schema compilation, caching, and detailed error reporting for validating
API responses and requests against JSON Schema definitions.

The manager integrates with SchemaLoader to load versioned schemas and provides both
generic and domain-specific validation methods for common contract types.

Architecture:
    - Uses fastjsonschema for fast, compiled validation
    - Schema compilation caching for optimal performance
    - Structured error reporting with JSON path information
    - Async-first API design

Example usage:

    ```python
    from app.clients.testing import ContractManager, ContractValidationError

    manager = ContractManager()

    # Validate user statistics
    try:
        validated_data = await manager.validate_user_stats({
            "student_id": "550e8400-e29b-41d4-a716-446655440000",
            "total_attempts": 5,
            "passed_attempts": 3,
            "completion_rate": 60.0,
            "average_score": 75.5,
            "last_attempt_at": "2025-12-17T14:30:00Z"
        })
        print("✓ Contract valid")
    except ContractValidationError as e:
        print(f"✗ Validation failed: {e.message}")
        for error in e.errors:
            print(f"  [{error['path']}] {error['message']}")
    ```

Features:
    - **Fast validation**: Uses fastjsonschema compiled validators for high performance
    - **Schema caching**: Automatically caches compiled validators
    - **Versioned schemas**: Supports multiple schema versions (v1, v2, latest)
    - **Rich error reporting**: Detailed error messages with JSON paths
    - **Type-safe**: Full type hints for IDE support
    - **Async-first**: Designed for async/await patterns

Supported contracts:
    - `user_stats` - User statistics and progress data
    - `attempt_detail` - Detailed test attempt information
    - `attempts_list` - Collection of test attempts

See Also:
    - `SchemaLoader`: Loads JSON schemas from disk
    - `ContractValidationError`: Exception raised on validation failures
"""

import logging
from pathlib import Path
from typing import Any, Callable

import fastjsonschema

from .schema_loader import SchemaLoader

logger = logging.getLogger(__name__)


class ContractValidationError(Exception):
    """Exception raised when contract validation fails.

    This exception provides detailed information about validation failures,
    including structured error details with JSON paths, messages, and constraint
    information from the fastjsonschema validator.

    Attributes:
        message (str): Human-readable summary of the validation failure.
        errors (list[dict[str, Any]]): List of structured validation errors. Each error contains:
            - `path` (str): JSON pointer to the invalid field (e.g., "data.student_id")
            - `message` (str): Human-readable error description
            - `rule` (str): Validation rule that failed (e.g., "type", "required")

    Example:

        ```python
        try:
            await manager.validate_user_stats(invalid_data)
        except ContractValidationError as e:
            print(f"Validation failed: {e.message}")
            print(f"Total errors: {len(e.errors)}")

            for error in e.errors:
                print(f"  Path: {error['path']}")
                print(f"  Issue: {error['message']}")
        ```
    """

    def __init__(self, message: str, errors: list[dict[str, Any]] | None = None) -> None:
        """Initialize validation error.

        Args:
            message: Human-readable error description.
            errors: List of structured validation error details.
        """
        self.message = message
        self.errors = errors or []
        super().__init__(self.message)

    def __str__(self) -> str:
        """Return string representation with error count."""
        error_count = len(self.errors)
        return f"{self.message} ({error_count} error{'s' if error_count != 1 else ''})"


class ContractManager:
    """Manager for API contract validation using fastjsonschema.

    This class provides a high-level interface for validating data against JSON schemas
    using the fastjsonschema library. It handles schema loading, compilation caching,
    and provides detailed error reporting.

    The manager uses fastjsonschema for fast validation with compiled validators:
    - JSON Schema Draft-04, Draft-06, Draft-07 support
    - Extremely fast validation (compiled to Python code)
    - Format validation support
    - Custom format handlers

    Attributes:
        schema_loader (SchemaLoader): Loader for JSON schemas.
        cache_size (int): Maximum number of compiled validators to cache.
        enable_caching (bool): Whether to enable validator caching.

    Example:

        ```python
        # Initialize with default settings
        manager = ContractManager()

        # Or customize caching behavior
        manager = ContractManager(
            cache_size=200,
            enable_caching=True
        )

        # Validate various contracts
        user_data = await manager.validate_user_stats(stats_dict)
        attempt = await manager.validate_attempt_detail(attempt_dict)
        attempts = await manager.validate_attempts_list(attempts_dict)
        ```
    """

    def __init__(
        self,
        schemas_dir: Path | None = None,
        cache_size: int = 100,
        enable_caching: bool = True,
    ) -> None:
        """Initialize contract manager with fastjsonschema validator.

        Args:
            schemas_dir: Path to schemas directory. If None, uses default location
                at app/clients/testing/schemas/.
            cache_size: Maximum number of compiled validators to keep in cache.
                Older validators are evicted when this limit is exceeded.
            enable_caching: Whether to enable validator caching. Set to False
                if schemas change frequently during development.
        """
        self.schema_loader = SchemaLoader(schemas_dir)
        self.cache_size = cache_size
        self.enable_caching = enable_caching

        # Cache for compiled validators: cache_key -> compiled_validator
        self._validators_cache: dict[str, Callable[[Any], Any]] = {}

        logger.debug("ContractManager initialized with fastjsonschema backend")

    async def validate(
        self,
        data: dict[str, Any] | list[Any],
        contract_name: str,
        version: str = "latest",
    ) -> dict[str, Any] | list[Any]:
        """Validate data against a contract schema.

        This is the primary validation method. It loads the appropriate schema,
        compiles it with fastjsonschema (or retrieves from cache), and validates the data.
        If validation fails, detailed error information is collected and raised.

        Args:
            data: The data to validate. Can be a dict or list depending on schema.
            contract_name: Name of the contract (e.g., "user_stats", "attempt_detail").
            version: Schema version to use (e.g., "v1", "v2", or "latest" for highest).

        Returns:
            The validated data (unchanged if validation succeeds).

        Raises:
            ContractValidationError: If data does not conform to the schema.
                Contains detailed error information with paths and messages.

        Example:

            ```python
            # Validate with specific version
            try:
                validated = await manager.validate(
                    data={"student_id": "uuid-here", ...},
                    contract_name="user_stats",
                    version="v1"
                )
                print("✓ Data valid")
            except ContractValidationError as e:
                print(f"✗ {e.message}")
                for err in e.errors:
                    print(f"  {err['path']}: {err['message']}")
            ```
        """
        try:
            # Load schema from disk (uses SchemaLoader cache)
            schema = await self.schema_loader.load(contract_name, version)

            # Get or compile validator
            cache_key = f"{contract_name}:{version}"
            validator = self._get_validator(schema, cache_key)

            # Validate data with fastjsonschema
            validator(data)

            logger.debug(f"✓ Contract validation passed: {contract_name}:{version}")
            return data

        except fastjsonschema.JsonSchemaValueException as e:
            # fastjsonschema validation error
            errors = self._format_fastjsonschema_error(e)
            error_msg = self._create_error_message(contract_name, errors)

            logger.warning(
                f"✗ Contract validation failed: {contract_name}:{version}",
                extra={"contract": contract_name, "version": version, "error_count": len(errors)},
            )

            raise ContractValidationError(error_msg, errors) from None

        except fastjsonschema.JsonSchemaDefinitionException as e:
            # Schema definition error (invalid schema)
            logger.error(f"Invalid schema definition for {contract_name}:{version}: {e}")
            raise ContractValidationError(f"Invalid schema for {contract_name}: {e!s}") from e

        except ContractValidationError:
            # Re-raise validation errors as-is
            raise

        except Exception as e:
            # Wrap unexpected errors
            logger.error(
                f"Unexpected error during validation of {contract_name}:{version}",
                exc_info=True,
            )
            raise ContractValidationError(f"Validation error for {contract_name}: {e!s}") from e

    async def validate_user_stats(self, data: dict[str, Any]) -> dict[str, Any]:
        """Validate UserStats contract (v1).

        Validates user statistics data including attempt counts, scores, and progress metrics.
        This is a convenience method that calls `validate()` with contract="user_stats" and version="v1".

        Expected schema fields:
            - student_id (uuid): Unique student identifier
            - total_attempts (int): Total number of test attempts
            - passed_attempts (int): Number of passed attempts
            - completion_rate (float): Percentage of completed tests (0-100)
            - average_score (float): Average test score
            - last_attempt_at (datetime): ISO 8601 timestamp of last attempt

        Args:
            data: UserStats data to validate.

        Returns:
            Validated data (unchanged if valid).

        Raises:
            ContractValidationError: If validation fails.

        Example:

            ```python
            stats = {
                "student_id": "550e8400-e29b-41d4-a716-446655440000",
                "total_attempts": 5,
                "passed_attempts": 3,
                "completion_rate": 60.0,
                "average_score": 75.5,
                "last_attempt_at": "2025-12-17T14:30:00Z"
            }
            validated = await manager.validate_user_stats(stats)
            ```
        """
        result = await self.validate(data, "user_stats", "v1")
        return result  # type: ignore[return-value]

    async def validate_attempts_list(self, data: dict[str, Any]) -> dict[str, Any]:
        """Validate AttemptsList contract (v1).

        Validates a collection of test attempts with pagination metadata.
        This is a convenience method that calls `validate()` with contract="attempts_list" and version="v1".

        Expected schema structure:
            - attempts (array): List of attempt objects
            - total (int): Total number of attempts
            - page (int): Current page number
            - page_size (int): Items per page

        Args:
            data: AttemptsList data to validate.

        Returns:
            Validated data (unchanged if valid).

        Raises:
            ContractValidationError: If validation fails.

        Example:

            ```python
            attempts_data = {
                "attempts": [
                    {"attempt_id": "...", "score": 85, ...},
                    {"attempt_id": "...", "score": 92, ...}
                ],
                "total": 2,
                "page": 1,
                "page_size": 20
            }
            validated = await manager.validate_attempts_list(attempts_data)
            ```
        """
        result = await self.validate(data, "attempts_list", "v1")
        return result  # type: ignore[return-value]

    async def validate_attempt_detail(self, data: dict[str, Any]) -> dict[str, Any]:
        """Validate AttemptDetail contract (v1).

        Validates detailed information about a single test attempt including answers,
        score breakdown, and completion status.
        This is a convenience method that calls `validate()` with contract="attempt_detail" and version="v1".

        Expected schema fields:
            - attempt_id (uuid): Unique attempt identifier
            - student_id (uuid): Student who made the attempt
            - test_id (uuid): Test being attempted
            - score (float): Final score
            - completed (bool): Whether attempt is complete
            - started_at (datetime): When attempt started
            - completed_at (datetime?): When attempt finished (nullable)
            - answers (array): List of answer objects

        Args:
            data: AttemptDetail data to validate.

        Returns:
            Validated data (unchanged if valid).

        Raises:
            ContractValidationError: If validation fails.

        Example:

            ```python
            attempt = {
                "attempt_id": "550e8400-e29b-41d4-a716-446655440000",
                "student_id": "...",
                "test_id": "...",
                "score": 85.5,
                "completed": True,
                "started_at": "2025-12-17T14:30:00Z",
                "completed_at": "2025-12-17T15:00:00Z",
                "answers": [...]
            }
            validated = await manager.validate_attempt_detail(attempt)
            ```
        """
        result = await self.validate(data, "attempt_detail", "v1")
        return result  # type: ignore[return-value]

    def clear_cache(self) -> None:
        """Clear all cached compiled validators.

        Useful during development when schemas are being modified frequently.
        After clearing, validators will be recompiled on next use.

        Example:

            ```python
            # After updating a schema file
            manager.clear_cache()
            # Next validation will use updated schema
            ```
        """
        self._validators_cache.clear()
        logger.debug("Validator cache cleared")

    def _get_validator(self, schema: dict[str, Any], cache_key: str) -> Callable[[Any], Any]:
        """Get compiled fastjsonschema validator from cache or compile new one.

        Args:
            schema: JSON Schema to compile.
            cache_key: Key to use for caching (typically "contract:version").

        Returns:
            Compiled fastjsonschema validator function.
        """
        # Return new validator if caching disabled
        if not self.enable_caching:
            return fastjsonschema.compile(schema)

        # Check cache
        if cache_key in self._validators_cache:
            logger.debug(f"Validator cache hit: {cache_key}")
            return self._validators_cache[cache_key]

        # Compile new validator
        validator = fastjsonschema.compile(schema)

        # Add to cache
        self._validators_cache[cache_key] = validator

        # Evict oldest if cache full
        if len(self._validators_cache) > self.cache_size:
            oldest_key = next(iter(self._validators_cache))
            del self._validators_cache[oldest_key]
            logger.debug(f"Evicted validator from cache: {oldest_key}")

        logger.debug(f"Compiled and cached validator: {cache_key}")
        return validator

    def _format_fastjsonschema_error(self, error: fastjsonschema.JsonSchemaValueException) -> list[dict[str, Any]]:
        """Format fastjsonschema validation error into structured list.

        Converts fastjsonschema exception into a consistent dictionary format
        with path, message, and rule for easier consumption.

        Args:
            error: JsonSchemaValueException from fastjsonschema.

        Returns:
            List of formatted error dictionaries with keys:
                - path: JSON path to invalid field
                - message: Human-readable error message
                - rule: Validation rule that failed
                - value: The invalid value (if available)
        """
        # fastjsonschema provides: message, value, name, path, rule, rule_definition
        path = ".".join(str(p) for p in error.path) if error.path else "root"

        # Extract rule from error name (e.g., "data.student_id must be string")
        rule = getattr(error, "rule", "unknown")
        if rule == "unknown":
            # Try to extract rule from message
            msg = str(error.message)
            if "must be" in msg:
                rule = "type"
            elif "required" in msg.lower():
                rule = "required"
            elif "minimum" in msg or "maximum" in msg:
                rule = "range"
            elif "pattern" in msg:
                rule = "pattern"
            elif "format" in msg:
                rule = "format"

        return [
            {
                "path": path,
                "message": str(error.message),
                "rule": rule,
                "value": error.value if hasattr(error, "value") else None,
            }
        ]

    def _create_error_message(self, contract_name: str, errors: list[dict[str, Any]]) -> str:
        """Create human-readable error message from validation errors.

        Formats the first 5 errors in a readable way with JSON paths and messages.
        Indicates if there are more errors beyond the displayed ones.

        Args:
            contract_name: Name of contract being validated.
            errors: List of formatted validation errors.

        Returns:
            Multi-line formatted error message.
        """
        if not errors:
            return f"Validation failed for contract '{contract_name}'"

        lines = [f"Contract '{contract_name}' validation failed:"]

        for i, error in enumerate(errors[:5], 1):
            path = error["path"]
            message = error["message"]
            rule = error.get("rule", "unknown")
            lines.append(f"  {i}. [{path}] {message} (rule: {rule})")

        if len(errors) > 5:
            lines.append(f"  ... and {len(errors) - 5} more error(s)")

        return "\n".join(lines)
