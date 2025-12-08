"""Frontend pages router with Jinja2 templates."""
from fastapi import APIRouter, Request
from fastapi.responses import HTMLResponse
from fastapi.templating import Jinja2Templates

from app.config import get_settings

settings = get_settings()
templates = Jinja2Templates(directory="templates")

router = APIRouter(tags=["Pages"])


def get_keycloak_urls() -> dict:
    """Generate Keycloak URLs for templates."""
    base_url = settings.KEYCLOAK_PUBLIC_URL
    realm = settings.KEYCLOAK_REALM
    client_id = settings.KEYCLOAK_CLIENT_ID
    
    return {
        "keycloak_account_url": f"{base_url}/realms/{realm}/account",
        "keycloak_register_url": f"{base_url}/realms/{realm}/protocol/openid-connect/registrations?client_id={client_id}&response_type=code&scope=openid&redirect_uri={settings.KEYCLOAK_REDIRECT_URI}",
    }


@router.get("/", response_class=HTMLResponse)
async def dashboard_page(request: Request):
    """Render dashboard page."""
    return templates.TemplateResponse(
        "dashboard.html",
        {
            "request": request,
            "active_page": "dashboard",
            **get_keycloak_urls()
        }
    )


@router.get("/login", response_class=HTMLResponse)
async def login_page(request: Request):
    """Render login page."""
    return templates.TemplateResponse(
        "login.html",
        {
            "request": request,
            **get_keycloak_urls()
        }
    )


@router.get("/profile", response_class=HTMLResponse)
async def profile_page(request: Request):
    """Render profile page."""
    return templates.TemplateResponse(
        "profile.html",
        {
            "request": request,
            "active_page": "profile",
            **get_keycloak_urls()
        }
    )


@router.get("/certificates", response_class=HTMLResponse)
async def certificates_page(request: Request):
    """Render certificates page."""
    return templates.TemplateResponse(
        "certificates.html",
        {
            "request": request,
            "active_page": "certificates",
            **get_keycloak_urls()
        }
    )


@router.get("/visits", response_class=HTMLResponse)
async def visits_page(request: Request):
    """Render visits page."""
    return templates.TemplateResponse(
        "visits.html",
        {
            "request": request,
            "active_page": "visits",
            **get_keycloak_urls()
        }
    )


@router.get("/register", response_class=HTMLResponse)
async def register_page(request: Request):
    """Render registration page."""
    return templates.TemplateResponse(
        "register.html",
        {
            "request": request,
            **get_keycloak_urls()
        }
    )


@router.get("/callback", response_class=HTMLResponse)
async def callback_page(request: Request):
    """Render OAuth callback page."""
    return templates.TemplateResponse(
        "callback.html",
        {
            "request": request
        }
    )
