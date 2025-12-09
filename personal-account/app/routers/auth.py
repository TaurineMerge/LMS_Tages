"""Authentication API endpoints."""
from fastapi import APIRouter, Request, Depends, status, Body, HTTPException
from fastapi.responses import RedirectResponse
from pydantic import BaseModel

from app.services.auth import auth_service, get_current_user
from app.schemas.common import (
    api_response, 
    token_response, 
    user_info_response, 
    message_response,
    register_request,
    register_response
)

router = APIRouter(prefix="/auth", tags=["Authentication"])
import logging
logger = logging.getLogger(__name__)


class refresh_token_request(BaseModel):
    """Refresh token request body."""
    refresh_token: str


class logout_request(BaseModel):
    """Logout request body."""
    refresh_token: str


class password_token_request(BaseModel):
    """Password token request body."""
    username: str
    password: str


@router.post(
    "/token",
    response_model=api_response[token_response],
    summary="Get access token",
    description="Get access token by username and password (OAuth2 password flow)."
)
async def get_token(body: password_token_request):
    """Get access token using username and password."""
    token = await auth_service.get_token_by_password(body.username, body.password)
    token_data = token_response(
        access_token=token.get("access_token", ""),
        refresh_token=token.get("refresh_token"),
        token_type=token.get("token_type", "Bearer"),
        expires_in=token.get("expires_in", 0),
        refresh_expires_in=token.get("refresh_expires_in"),
        scope=token.get("scope")
    )
    return api_response(data=token_data, message="Token obtained successfully")


@router.get(
    "/login",
    summary="Redirect to Keycloak login",
    description="Redirects the user to the Keycloak login page for the Student realm."
)
async def login():
    """Initiate login flow."""
    login_url = auth_service.get_login_url()
    # Log the redirect URL to help diagnose 'client not found' errors
    try:
        logger.debug("Redirecting to Keycloak login URL: %s", login_url)
    except Exception:
        pass
    return RedirectResponse(url=login_url)


@router.get(
    "/callback",
    response_model=api_response[token_response],
    summary="Handle Keycloak callback",
    description="Exchanges the authorization code for an access token."
)
async def callback(code: str):
    """Handle login callback."""
    logger.info(f"Callback received with code: {code[:20]}..." if len(code) > 20 else f"Callback received with code: {code}")
    try:
        token = await auth_service.exchange_code_for_token(code)
        token_data = token_response(
            access_token=token.get("access_token", ""),
            refresh_token=token.get("refresh_token"),
            token_type=token.get("token_type", "Bearer"),
            expires_in=token.get("expires_in", 0),
            refresh_expires_in=token.get("refresh_expires_in"),
            scope=token.get("scope")
        )
        logger.info("Callback successful, returning token")
        return api_response(data=token_data, message="Login successful")
    except HTTPException as e:
        logger.error(f"Callback failed: {e.detail}")
        raise
    except Exception as e:
        logger.error(f"Unexpected callback error: {str(e)}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Callback failed: {str(e)}"
        )


@router.post(
    "/refresh",
    response_model=api_response[token_response],
    summary="Refresh access token",
    description="Exchanges a refresh token for a new access token."
)
async def refresh(body: refresh_token_request):
    """Refresh access token."""
    token = await auth_service.refresh_token(body.refresh_token)
    token_data = token_response(
        access_token=token.get("access_token", ""),
        refresh_token=token.get("refresh_token"),
        token_type=token.get("token_type", "Bearer"),
        expires_in=token.get("expires_in", 0),
        refresh_expires_in=token.get("refresh_expires_in"),
        scope=token.get("scope")
    )
    return api_response(data=token_data, message="Token refreshed successfully")


@router.post(
    "/logout",
    response_model=api_response[message_response],
    summary="Logout user",
    description="Invalidates the user's session by revoking the refresh token."
)
async def logout(body: logout_request):
    """Logout user."""
    await auth_service.logout(body.refresh_token)
    return api_response(data=message_response(message="Logged out successfully"))


@router.get(
    "/me",
    response_model=api_response[user_info_response],
    summary="Get current user info",
    description="Returns information about the currently authenticated user."
)
async def get_me(user: dict = Depends(get_current_user)):
    """Get current user profile."""
    user_data = user_info_response(
        sub=user.get("sub", ""),
        email=user.get("email"),
        preferred_username=user.get("preferred_username"),
        name=user.get("name"),
        given_name=user.get("given_name"),
        family_name=user.get("family_name"),
        email_verified=user.get("email_verified")
    )
    return api_response(data=user_data)


@router.get(
    "/register",
    summary="Redirect to Keycloak registration",
    description="Redirects the user to the Keycloak registration page."
)
async def register_redirect():
    """Redirect to Keycloak registration page."""
    register_url = auth_service.get_register_url()
    return RedirectResponse(url=register_url)


@router.post(
    "/register",
    response_model=api_response[register_response],
    status_code=status.HTTP_201_CREATED,
    summary="Register new user",
    description="Creates a new user account in Keycloak. Available for both students and teachers."
)
async def register(body: register_request):
    """Register a new user via API."""
    result = await auth_service.register_user(
        username=body.username,
        email=body.email,
        password=body.password,
        first_name=body.first_name,
        last_name=body.last_name
    )
    
    response_data = register_response(
        user_id=result["user_id"],
        username=result["username"],
        email=result["email"]
    )
    return api_response(data=response_data, message="User registered successfully")


