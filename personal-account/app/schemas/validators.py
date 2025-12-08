"""Utility validators shared across Pydantic schemas."""
from __future__ import annotations

import re
from typing import Any

_SQL_META_CHARS = re.compile(r"(--|;|/\*|\*/)")
_SQL_UNSAFE_PATTERNS = [
    re.compile(r"(?i)\bunion\b\s+\bselect\b"),
    re.compile(r"(?i)\bord\b\s+1\s*=\s*1"),
    re.compile(r"(?i)\bexec\b"),
]


def ensure_safe_string(value: str | None, field_name: str) -> str | None:
    """Ensure that the provided string doesn't contain obvious SQL-injection payloads."""
    if value is None:
        return None

    if _SQL_META_CHARS.search(value):
        raise ValueError(f"Поле '{field_name}' содержит запрещённые символы SQL")

    for pattern in _SQL_UNSAFE_PATTERNS:
        if pattern.search(value):
            raise ValueError(
                f"Поле '{field_name}' содержит потенциально опасную SQL-конструкцию",
            )

    return value


def ensure_safe_mapping(mapping: dict[str, Any] | None, field_name: str) -> dict[str, Any] | None:
    """Recursively validate every string inside a mapping."""
    if not mapping:
        return mapping

    validated: dict[str, Any] = {}
    for key, val in mapping.items():
        if isinstance(val, str):
            validated[key] = ensure_safe_string(val, f"{field_name}.{key}")
        elif isinstance(val, dict):
            validated[key] = ensure_safe_mapping(val, f"{field_name}.{key}")
        else:
            validated[key] = val
    return validated
