"""Frontend pages router with Jinja2 templates."""

import logging
from uuid import UUID

from fastapi import APIRouter, Request
from fastapi.responses import (
    HTMLResponse,
    RedirectResponse,
)  # <--- Добавили RedirectResponse
from fastapi.templating import Jinja2Templates
from opentelemetry import trace

from app.config import get_settings
from app.core.security import JWTValidator
from app.services import stats_service, stats_worker
from app.telemetry import traced

settings = get_settings()
templates = Jinja2Templates(directory="templates")

router = APIRouter(tags=["Pages"])

logger = logging.getLogger(__name__)

# JWT validator for token decoding
jwt_validator = JWTValidator(
    keycloak_server_url=settings.KEYCLOAK_PUBLIC_URL,  # Используем PUBLIC_URL для issuer
    realm=settings.KEYCLOAK_REALM,
    client_id=settings.KEYCLOAK_CLIENT_ID,
)


def _render_template_safe(template_name: str, context: dict):
    """Render a Jinja2 template and attach telemetry/logging on exception.

    This helper records exceptions on the current span and logs the error
    with trace/span identifiers so failures in template rendering are visible
    in logs and traces.
    """
    request = context.get("request")
    if not request:
        logger.error("Request object missing in template context for %s", template_name)
        raise ValueError("Request object is required for TemplateResponse")

    try:
        # TemplateResponse will be returned to FastAPI and rendered by Starlette
        context["prefix"] = settings.url_prefix  # Динамический prefix из settings

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


@router.get("/", response_class=HTMLResponse)
@traced("pages.root")
async def root_page(request: Request):
    """
    Главная страница - доступна как /account и /account/.
    Логика:
    1. Если у пользователя есть кука 'access_token' -> редирект в /dashboard
    2. Иначе -> показываем Лендинг
    """
    token = request.cookies.get("access_token")

    if token:
        return RedirectResponse(url=f"{settings.url_prefix}/dashboard")

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


@router.get("/statistics", response_class=HTMLResponse)
@traced("pages.statistics")
async def statistics_page(request: Request):
    """Render statistics page.

    This page displays aggregated user statistics loaded from Redis cache.
    """
    # Get token from cookie
    access_token = request.cookies.get("access_token")
    if not access_token:
        # Redirect to login if not authenticated
        return RedirectResponse(url="/api/v1/auth/login", status_code=302)

    # Decode token to get user info
    try:
        token_payload = await jwt_validator.validate_token(access_token)
        student_id = token_payload.sub
    except Exception:
        # Token invalid, redirect to login
        return RedirectResponse(url="/api/v1/auth/login", status_code=302)

    # Fetch stats from backend
    try:
        # logging.debug("Fetching stats for user %s", student_id)
        # await stats_worker.fetch_from_testing(UUID(student_id))
        # logging.debug("Processing raws for user %s", student_id)
        # await stats_worker.process_raws(UUID(student_id))
        logging.debug("Getting stats for user %s", student_id)
        stats = await stats_service.get_user_statistics(UUID(student_id))
        logging.debug("Fetched stats for user %s: %s", student_id, stats)
    except Exception as e:
        logger.error(f"Failed to fetch stats for user {student_id}: {e}")
        stats = {}

    return _render_template_safe(
        "statistics.hbs",
        {
            "request": request,
            "active_page": "statistics",
            "stats": stats,
            "student_id": student_id,
            **get_keycloak_urls(),
        },
    )


@router.post("/statistics/refresh")
@traced("pages.statistics_refresh")
async def statistics_refresh(request: Request):
    """Refresh statistics and redirect back."""
    access_token = request.cookies.get("access_token")
    if not access_token:
        return RedirectResponse(url="/api/v1/auth/login", status_code=302)

    try:
        token_payload = await jwt_validator.validate_token(access_token)
        student_id = UUID(token_payload.sub)

        # 1. Fetch fresh data from testing service
        await stats_worker.fetch_for_student(student_id)

        # 2. Process raw data into business tables
        # Note: In production, this might be too slow for a request-response cycle
        # but for this practice it's fine.
        await stats_worker.processor.process_raw_user_stats()
        await stats_worker.processor.process_raw_attempts()

        # 3. Recalculate aggregated stats
        await stats_service.refresh_user_statistics(student_id)

        logger.info("Successfully refreshed stats for student %s", student_id)
    except Exception as e:
        logger.error(f"Failed to refresh stats: {e}")

    return RedirectResponse(url=f"{settings.url_prefix}/statistics", status_code=303)


@router.get("/register", response_class=HTMLResponse)
@traced("pages.register")
async def register_page(request: Request):
    """Render registration page."""
    return _render_template_safe("register.hbs", {"request": request, **get_keycloak_urls()})


@router.get("/callback", response_class=HTMLResponse)
async def callback_page(request: Request):
    """Render OAuth callback page."""
    return _render_template_safe("callback.hbs", {"request": request})
