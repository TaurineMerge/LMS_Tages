"""Authentication API endpoints."""
from fastapi import APIRouter, Request, Depends, status, Body, HTTPException
from fastapi.responses import RedirectResponse, JSONResponse
from typing import Optional
from pydantic import BaseModel
from app.telemetry import traced
from app.services.auth import auth_service
from app.core.security import get_current_user, TokenPayload
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


class Refresh_Token_Request(BaseModel):
    """Refresh token request body."""
    refresh_token: str


class Logout_Request(BaseModel):
    """Logout request body."""
    refresh_token: str = None  # Optional for cookie-based auth


class Password_Token_Request(BaseModel):
    """Password token request body."""
    username: str
    password: str


@router.post(
    "/token",
    response_model=api_response[token_response],
    summary="Get access token",
    description="Get access token by username and password (OAuth2 password flow)."
)
@traced("router.auth.token")
async def get_token(body: Password_Token_Request):
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
@traced("router.auth.login", record_result=True)
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
    summary="Handle Keycloak callback",
    description="Exchanges the authorization code for an access token and sets secure cookies."
)
@traced("router.auth.callback")
async def callback(code: str):
    """Handle login callback and set secure cookie with access token."""
    logger.info(f"Callback received with code: {code[:20]}..." if len(code) > 20 else f"Callback received with code: {code}")
    try:
        token = await auth_service.exchange_code_for_token(code)

        # Create response with JSON data
        token_data = token_response(
            access_token=token.get("access_token", ""),
            refresh_token=token.get("refresh_token"),
            token_type=token.get("token_type", "Bearer"),
            expires_in=token.get("expires_in", 0),
            refresh_expires_in=token.get("refresh_expires_in"),
            scope=token.get("scope")
        )

        response_data = api_response(data=token_data, message="Login successful")
        response = JSONResponse(content=response_data.model_dump())

        # Set secure HTTP-only cookie for access token
        response.set_cookie(
            key="access_token",
            value=token.get("access_token", ""),
            max_age=token.get("expires_in", 3600),
            httponly=True,
            secure=True,  # Requires HTTPS in production
            samesite="Lax",
            path="/"
        )

        # Set refresh token cookie if provided
        if token.get("refresh_token"):
            response.set_cookie(
                key="refresh_token",
                value=token.get("refresh_token", ""),
                max_age=token.get("refresh_expires_in", 86400),
                httponly=True,
                secure=True,  # Requires HTTPS in production
                samesite="Lax",
                path="/"
            )

        logger.info("Callback successful, tokens set in secure cookies")
        return response

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
@traced("router.auth.refresh")
async def refresh(body: Refresh_Token_Request):
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
    description="Logs out the user by clearing server-side session and cookies."
)
@traced("router.auth.logout")
async def logout(body: Optional[Logout_Request] = None):
    """Logout user.
    
    Works with both:
    1. Traditional: body with refresh_token for logout from other clients
    2. Cookie-based: empty body, relies on server-side session invalidation
    """
    if body and body.refresh_token:
        await auth_service.logout(body.refresh_token)
    
    # Create response that clears cookies
    response = JSONResponse(
        content=api_response(data=message_response(message="Logged out successfully")).model_dump()
    )
    
    # Clear access_token cookie
    response.delete_cookie(
        key="access_token",
        path="/"
    )
    
    # Clear refresh_token cookie
    response.delete_cookie(
        key="refresh_token",
        path="/"
    )
    
    return response





@router.get(
    "/me",
    response_model=api_response[user_info_response],
    summary="Get current user info",
    description="Returns information about the currently authenticated user."
)
@traced("router.auth.me")
async def get_me(user: TokenPayload = Depends(get_current_user)):
    """Get current user profile."""
    user_data = user_info_response(
        sub=user.sub,
        email=user.email,
        preferred_username=user.preferred_username,
        name=user.name,
        given_name=user.given_name,
        family_name=user.family_name,
        email_verified=user.email_verified
    )
    return api_response(data=user_data)





@router.get(
    "/register",
    summary="Redirect to Keycloak registration",
    description="Redirects the user to the Keycloak registration page."
)
@traced("router.auth.register_redirect")
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
@traced("router.auth.register")
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





