"""Telemetry helpers for tracing application code."""
from __future__ import annotations

import inspect
import json
import logging
from functools import wraps
from typing import Any, Callable, ParamSpec, TypeVar

from opentelemetry import trace
from opentelemetry.trace import StatusCode

P = ParamSpec("P")
T = TypeVar("T")

_tracer = trace.get_tracer("personal-account")
_logger = logging.getLogger(__name__)

# Maximum length for serialized attribute values to avoid overloading trace backend
_MAX_ATTR_LEN = 1024


def _safe_serialize(value: Any, max_len: int = _MAX_ATTR_LEN) -> str:
    """Safely serialize a value to string for span attributes."""
    try:
        if value is None:
            return "None"
        if isinstance(value, (str, int, float, bool)):
            s = str(value)
        elif isinstance(value, (list, tuple)):
            s = json.dumps(value[:10], default=str)  # limit list preview
        elif isinstance(value, dict):
            # Limit dict preview to first 10 keys
            preview = {k: v for i, (k, v) in enumerate(value.items()) if i < 10}
            s = json.dumps(preview, default=str)
        elif hasattr(value, "__dict__"):
            # Pydantic models or dataclasses
            d = getattr(value, "model_dump", getattr(value, "dict", lambda: value.__dict__))
            if callable(d):
                d = d()
            preview = {k: v for i, (k, v) in enumerate(d.items()) if i < 10}
            s = json.dumps(preview, default=str)
        else:
            s = repr(value)
        return s[:max_len] + ("..." if len(s) > max_len else "")
    except Exception:
        return "<unserializable>"


def _format_args(func: Callable[..., Any], args: tuple, kwargs: dict) -> dict[str, str]:
    """Format function arguments as a dict of serialized values."""
    sig = inspect.signature(func)
    params = list(sig.parameters.keys())
    result: dict[str, str] = {}
    for i, arg in enumerate(args):
        key = params[i] if i < len(params) else f"arg{i}"
        # Skip 'self' / 'cls'
        if key in ("self", "cls"):
            continue
        result[f"arg.{key}"] = _safe_serialize(arg)
    for k, v in kwargs.items():
        result[f"arg.{k}"] = _safe_serialize(v)
    return result


def traced(
    span_name: str | None = None,
    *,
    record_args: bool = True,
    record_result: bool = True,
) -> Callable[[Callable[P, T]], Callable[P, T]]:
    """
    Decorator for adding OpenTelemetry spans around sync/async callables.

    Parameters
    ----------
    span_name : str | None
        Custom span name. Defaults to function qualified name.
    record_args : bool
        Whether to record function arguments as span attributes.
    record_result : bool
        Whether to record function return value as span attribute.
    """

    def decorator(func: Callable[P, T]) -> Callable[P, T]:
        name = span_name or func.__qualname__

        if inspect.iscoroutinefunction(func):

            @wraps(func)
            async def async_wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
                with _tracer.start_as_current_span(name) as span:
                    _attach_code_attributes(span, func)
                    if record_args:
                        _attach_args(span, func, args, kwargs)
                    try:
                        result = await func(*args, **kwargs)
                        if record_result:
                            _attach_result(span, result)
                        span.set_status(StatusCode.OK)
                        return result
                    except Exception as exc:
                        _attach_exception(span, exc)
                        raise

            return async_wrapper

        @wraps(func)
        def sync_wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
            with _tracer.start_as_current_span(name) as span:
                _attach_code_attributes(span, func)
                if record_args:
                    _attach_args(span, func, args, kwargs)
                try:
                    result = func(*args, **kwargs)
                    if record_result:
                        _attach_result(span, result)
                    span.set_status(StatusCode.OK)
                    return result
                except Exception as exc:
                    _attach_exception(span, exc)
                    raise

        return sync_wrapper

    return decorator


def get_tracer() -> trace.Tracer:
    """Expose module tracer for advanced scenarios."""
    return _tracer


def _attach_code_attributes(span: trace.Span, func: Callable[..., Any]) -> None:
    """Attach code location attributes to span."""
    if not span.is_recording():
        return
    span.set_attribute("code.namespace", func.__module__)
    span.set_attribute("code.function", func.__qualname__)
    # Try to attach line number
    try:
        span.set_attribute("code.lineno", inspect.getsourcelines(func)[1])
        span.set_attribute("code.filepath", inspect.getfile(func))
    except (OSError, TypeError):
        pass


def _attach_args(
    span: trace.Span, func: Callable[..., Any], args: tuple, kwargs: dict
) -> None:
    """Attach function arguments to span."""
    if not span.is_recording():
        return
    for k, v in _format_args(func, args, kwargs).items():
        span.set_attribute(k, v)


def _attach_result(span: trace.Span, result: Any) -> None:
    """Attach function return value to span."""
    if not span.is_recording():
        return
    span.set_attribute("result.type", type(result).__name__)
    span.set_attribute("result.value", _safe_serialize(result))
    # For collections, record count
    if hasattr(result, "__len__"):
        try:
            span.set_attribute("result.count", len(result))
        except TypeError:
            pass


def _attach_exception(span: trace.Span, exc: Exception) -> None:
    """Record exception details on span."""
    span.set_status(StatusCode.ERROR, str(exc))
    span.record_exception(exc)
    span.set_attribute("error.type", type(exc).__name__)
    span.set_attribute("error.message", str(exc)[:_MAX_ATTR_LEN])
