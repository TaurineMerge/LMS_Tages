"""Frontend pages router with Jinja2 templates."""

from fastapi import APIRouter, Request
from fastapi.responses import (
    HTMLResponse,
    RedirectResponse,
)  # <--- Добавили RedirectResponse
from fastapi.templating import Jinja2Templates

from app.config import get_settings
from app.telemetry import traced
from opentelemetry import trace
import logging

settings = get_settings()
templates = Jinja2Templates(directory="templates")

router = APIRouter(tags=["Pages"])

logger = logging.getLogger(__name__)


def _render_template_safe(template_name: str, context: dict):
    """Render a Jinja2 template and attach telemetry/logging on exception.

    This helper records exceptions on the current span and logs the error
    with trace/span identifiers so failures in template rendering are visible
    in logs and traces.
    """
    try:
        # TemplateResponse will be returned to FastAPI and rendered by Starlette
        resp = templates.TemplateResponse(template_name, context)
        # Add trace headers to response (best-effort)
        try:
            span = trace.get_current_span()
            span_ctx = span.get_span_context()
            if span_ctx and span_ctx.is_valid:
                resp.headers["X-Trace-Id"] = format(span_ctx.trace_id, "032x")
                resp.headers["X-Span-Id"] = format(span_ctx.span_id, "016x")
        except Exception:
            # Don't fail rendering if tracing helpers fail
            logger.debug("Failed to attach trace headers to template response")
        return resp
    except Exception as exc:
        # Attach exception to current span and log with trace identifiers
        span = trace.get_current_span()
        try:
            span.record_exception(exc)
            span.set_attribute("error.rendering_template", template_name)
            span.set_status(trace.StatusCode.ERROR)
        except Exception:
            # ignore telemetry failures
            pass
        # Log with trace/span ids
        try:
            span_ctx = span.get_span_context()
            trace_id = format(span_ctx.trace_id, "032x") if span_ctx and span_ctx.is_valid else "-"
            span_id = format(span_ctx.span_id, "016x") if span_ctx and span_ctx.is_valid else "-"
        except Exception:
            trace_id = "-"
            span_id = "-"
        logger.exception(
            "Template render failed (%s) trace_id=%s span_id=%s",
            template_name,
            trace_id,
            span_id,
        )
        raise


def get_keycloak_urls() -> dict:
    """Generate Keycloak URLs for templates."""
    base_url = settings.KEYCLOAK_PUBLIC_URL
    realm = settings.KEYCLOAK_REALM
    client_id = settings.KEYCLOAK_CLIENT_ID
    account_url = f"{base_url}/realms/{realm}/account"
    register_url = (
        f"{base_url}/realms/{realm}/protocol/openid-connect/registrations?client_id={client_id}"
        f"&response_type=code&scope=openid&redirect_uri={settings.KEYCLOAK_REDIRECT_URI}"
    )
    # Debug log to help diagnose 'client not found' issues from Keycloak
    logger.debug(
        "Keycloak URLs: account=%s register=%s client_id=%s",
        account_url,
        register_url,
        client_id,
    )

    return {
        "keycloak_account_url": account_url,
        "keycloak_register_url": register_url,
    }


@router.get("/account", response_class=HTMLResponse)
@traced("pages.root")
async def root_page(request: Request):
    """
    Главная страница (Root).
    Логика:
    1. Если у пользователя есть кука 'access_token' -> редирект в /dashboard
    2. Иначе -> показываем Лендинг
    """
    token = request.cookies.get("access_token")

    if token:
        return RedirectResponse(url="/dashboard")  # Full path with /account prefix

    return _render_template_safe("index.hbs", {"request": request, **get_keycloak_urls()})


@router.get("/dashboard", response_class=HTMLResponse)  # <-- Теперь дэшборд здесь
@traced("pages.dashboard")
async def dashboard_page(request: Request):
    """Render dashboard page (Protected Area)."""

    # (Опционально) Можно добавить проверку: если нет куки, редирект на /
    # token = request.cookies.get("access_token")
    # if not token:
    #     return RedirectResponse(url="/")

    return _render_template_safe(
        "dashboard.hbs",
        {"request": request, "active_page": "dashboard", **get_keycloak_urls()},
    )


@router.get("/profile", response_class=HTMLResponse)
@traced("pages.profile")
async def profile_page(request: Request):
    """Render profile page."""
    return _render_template_safe(
        "profile.hbs",
        {"request": request, "active_page": "profile", **get_keycloak_urls()},
    )


@router.get("/certificates", response_class=HTMLResponse)
@traced("pages.certificates")
async def certificates_page(request: Request):
    """Render certificates page."""
    return _render_template_safe(
        "certificates.hbs",
        {"request": request, "active_page": "certificates", **get_keycloak_urls()},
    )


@router.get("/visits", response_class=HTMLResponse)
@traced("pages.visits")
async def visits_page(request: Request):
    """Render visits page."""
    return _render_template_safe(
        "visits.hbs",
        {"request": request, "active_page": "visits", **get_keycloak_urls()},
    )


@router.get("/register", response_class=HTMLResponse)
@traced("pages.register")
async def register_page(request: Request):
    """Render registration page."""
    return _render_template_safe("register.hbs", {"request": request, **get_keycloak_urls()})


@router.get("/callback", response_class=HTMLResponse)
async def callback_page(request: Request):
    """Render OAuth callback page."""
    return _render_template_safe("callback.hbs", {"request": request})
