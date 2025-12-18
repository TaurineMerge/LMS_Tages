"""Schema loader for JSON contract validation.

This module provides asynchronous loading and caching of JSON Schema files used for
validating API contracts. It supports versioned schema management, automatic version
discovery, and intelligent caching for optimal performance.

The SchemaLoader is designed to work seamlessly with the ContractManager, handling
all file I/O and schema parsing while the ContractManager focuses on validation logic.

Architecture:
    - **Async I/O**: Non-blocking file operations using aiofiles
    - **Smart Caching**: In-memory cache with per-schema:version keys
    - **Version Management**: Automatic "latest" version resolution
    - **Validation**: Basic JSON Schema structure validation

Directory structure:
    ```
    schemas/
    ├── user_stats/
    │   ├── v1.json
    │   └── v2.json
    ├── attempt_detail/
    │   └── v1.json
    └── attempts_list/
        └── v1.json
    ```

Example usage:

    ```python
    from app.clients.testing.schema_loader import SchemaLoader

    # Initialize loader
    loader = SchemaLoader()

    # Load specific version
    schema_v1 = await loader.load("attempt_detail", "v1")

    # Load latest version automatically
    latest_schema = await loader.load("user_stats", "latest")

    # Clear cache if schemas change
    loader.clear_cache()
    ```

Features:
    - **Async I/O**: Non-blocking file operations for high concurrency
    - **Intelligent caching**: Reduces disk I/O for frequently used schemas
    - **Version discovery**: Automatically finds latest version when requested
    - **Structure validation**: Ensures loaded schemas are valid JSON Schema
    - **Error handling**: Clear error messages for missing or invalid schemas

See Also:
    - `ContractManager`: Uses SchemaLoader to validate contracts
    - JSON Schema specification: https://json-schema.org/
"""

import json
import logging
from pathlib import Path
from typing import Any

import aiofiles

logger = logging.getLogger(__name__)


class SchemaLoader:
    """Asynchronous loader for JSON schemas with caching support.

    This class handles loading and caching of JSON Schema files from disk.
    Schemas are organized in directories by contract name, with files named
    by version (e.g., v1.json, v2.json).

    Attributes:
        schemas_dir (Path): Root directory containing schema subdirectories.
        _schemas_cache (Dict): In-memory cache of loaded schemas.

    Example:

        ```python
        loader = SchemaLoader()

        # Load a specific version
        schema = await loader.load("user_stats", "v1")

        # Load latest version
        schema = await loader.load("attempt_detail", "latest")

        # Clear cache if needed
        loader.clear_cache()
        ```
    """

    def __init__(self, schemas_dir: Path | None = None):
        """Initialize schema loader.

        Args:
            schemas_dir: Path to schemas directory. If None, uses
                app/clients/testing/schemas relative to this module.
        """
        self.schemas_dir = schemas_dir or self._get_default_schemas_dir()
        self._ensure_schemas_dir()
        self._schemas_cache: dict[str, dict[str, Any]] = {}

    async def load(self, contract_name: str, version: str = "latest") -> dict[str, Any]:
        """Load a JSON schema for the specified contract and version.

        Schemas are cached after first load. The cache key is constructed from
        contract_name and version to allow different versions to be loaded
        independently.

        Args:
            contract_name: Name of the contract (e.g., "user_stats", "attempt_detail").
            version: Schema version (e.g., "v1", "v2") or "latest" to auto-select
                the highest available version.

        Returns:
            Dictionary containing the JSON Schema.

        Raises:
            FileNotFoundError: If contract directory or schema file does not exist.
            ValueError: If schema file contains invalid JSON or is malformed.

        Example:

            ```python
            # Load attempt detail schema
            schema = await loader.load("attempt_detail", "v1")

            # Automatically find and load latest version
            latest_schema = await loader.load("user_stats", "latest")
            ```
        """
        cache_key = f"{contract_name}:{version}"

        # Check cache first
        if cache_key in self._schemas_cache:
            logger.debug(f"Schema cache hit: {cache_key}")
            return self._schemas_cache[cache_key]

        # Find schema file
        schema_path = self._find_schema_file(contract_name, version)

        # Load schema from disk
        schema = await self._load_schema_file(schema_path)

        # Validate basic schema structure
        self._validate_schema_structure(schema, contract_name)

        # Cache the result
        self._schemas_cache[cache_key] = schema
        logger.debug(f"Schema loaded and cached: {cache_key} from {schema_path}")

        return schema

    def _find_schema_file(self, contract_name: str, version: str) -> Path:
        """Find the path to a schema file.

        Args:
            contract_name: Name of the contract directory.
            version: Specific version or "latest".

        Returns:
            Path to the schema file.

        Raises:
            FileNotFoundError: If contract or schema file not found.
        """
        contract_dir = self.schemas_dir / contract_name

        if not contract_dir.exists():
            raise FileNotFoundError(f"Contract directory not found: {contract_name}")

        if version == "latest":
            # Find latest version by sorting filenames
            json_files = sorted(contract_dir.glob("*.json"))
            if not json_files:
                raise FileNotFoundError(f"No schema files found for contract: {contract_name}")

            # Return the last file (highest version)
            schema_path = json_files[-1]
            logger.debug(f"Using latest schema for {contract_name}: {schema_path.name}")
            return schema_path

        # Look for specific version
        schema_path = contract_dir / f"{version}.json"
        if not schema_path.exists():
            raise FileNotFoundError(f"Schema not found: {contract_name}:{version} at {schema_path}")

        return schema_path

    async def _load_schema_file(self, schema_path: Path) -> dict[str, Any]:
        """Load and parse a JSON schema file asynchronously.

        Args:
            schema_path: Path to the schema file.

        Returns:
            Parsed JSON schema as dictionary.

        Raises:
            ValueError: If JSON is malformed.
            IOError: If file cannot be read.
        """
        try:
            async with aiofiles.open(schema_path, encoding="utf-8") as f:
                content = await f.read()
                return json.loads(content)
        except json.JSONDecodeError as e:
            logger.error(f"Invalid JSON in schema file {schema_path}: {e}")
            raise ValueError(f"Invalid JSON schema file: {schema_path}") from e
        except OSError as e:
            logger.error(f"Error loading schema file {schema_path}: {e}")
            raise

    def _validate_schema_structure(self, schema: dict[str, Any], contract_name: str) -> None:
        """Validate that schema has required structure.

        Args:
            schema: The schema dictionary to validate.
            contract_name: Name of contract (for logging).

        Raises:
            ValueError: If schema is missing required properties.
        """
        if not isinstance(schema, dict):
            raise ValueError(f"Schema must be a JSON object, got {type(schema)}")

        if "$schema" not in schema:
            logger.warning(f"Schema {contract_name} missing $schema property (recommended)")

        if "type" not in schema:
            raise ValueError(f"Schema {contract_name} must have 'type' property")

    def _get_default_schemas_dir(self) -> Path:
        """Get default schemas directory path.

        Returns:
            Path to app/clients/testing/schemas directory.
        """
        current_dir = Path(__file__).parent
        return current_dir / "schemas"

    def _ensure_schemas_dir(self) -> None:
        """Create schemas directory if it does not exist."""
        self.schemas_dir.mkdir(parents=True, exist_ok=True)
        logger.debug(f"Schemas directory ready: {self.schemas_dir}")

    def clear_cache(self) -> None:
        """Clear the schema cache.

        Call this if schema files are modified on disk and you need to
        reload them without restarting the application.
        """
        self._schemas_cache.clear()
        logger.debug("Schema cache cleared")
