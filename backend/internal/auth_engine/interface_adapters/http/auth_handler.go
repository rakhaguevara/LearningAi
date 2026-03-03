package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/adaptive-ai-learn/backend/internal/auth_engine/domain"
	"github.com/adaptive-ai-learn/backend/internal/auth_engine/usecase"
)

type AuthHandler struct {
	registerUC  *usecase.RegisterUseCase
	loginUC     *usecase.LoginUseCase
	googleUC    *usecase.GoogleOAuthUseCase
	refreshUC   *usecase.RefreshTokenUseCase
	logoutUC    *usecase.LogoutUseCase
	meUC        *usecase.MeUseCase
	frontendURL string
}

func NewAuthHandler(
	registerUC *usecase.RegisterUseCase,
	loginUC *usecase.LoginUseCase,
	googleUC *usecase.GoogleOAuthUseCase,
	refreshUC *usecase.RefreshTokenUseCase,
	logoutUC *usecase.LogoutUseCase,
	meUC *usecase.MeUseCase,
	frontendURL string,
) *AuthHandler {
	return &AuthHandler{
		registerUC:  registerUC,
		loginUC:     loginUC,
		googleUC:    googleUC,
		refreshUC:   refreshUC,
		logoutUC:    logoutUC,
		meUC:        meUC,
		frontendURL: frontendURL,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req usecase.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "code": "BAD_REQUEST"})
		return
	}

	if err := h.registerUC.Execute(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "REGISTER_FAILED"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req usecase.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "code": "BAD_REQUEST"})
		return
	}

	req.IP = c.ClientIP()
	req.Device = c.GetHeader("User-Agent")

	tokenPair, err := h.loginUC.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "UNAUTHORIZED"})
		return
	}

	h.setTokenCookies(c, tokenPair)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Generate a secure random state in prod, tying it to a session cookie.
	// For simplicity, hardcoding a state.
	state := "secure-state"
	url := h.googleUC.GetLoginURL(state)
	c.Redirect(http.StatusFound, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	// state := c.Query("state")
	// TODO: verify state matches the one stored in session

	tokenPair, err := h.googleUC.HandleCallback(c.Request.Context(), code, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		log.Printf("Google OAuth failed: %v", err)
		c.Redirect(http.StatusFound, h.frontendURL+"/login?error=oauth_failed")
		return
	}

	h.setTokenCookies(c, tokenPair)
	c.Redirect(http.StatusFound, h.frontendURL+"/dashboard")
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	// Let's extract the refresh token from httpOnly cookie
	rawToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token", "code": "UNAUTHORIZED"})
		return
	}

	// Assuming the UI also passes the UserID implicitly, or we decode JWT to grab sub which means empty strings for userID is okay depending on implementation
	// We will fake an empty userID string and let the UseCase rely on the JWT extracted ID.
	tokenPair, err := h.refreshUC.Execute(c.Request.Context(), "bypass", rawToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": "UNAUTHORIZED"})
		return
	}

	h.setTokenCookies(c, tokenPair)
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract userID from context (set by AuthMiddleware)
	userIDStr, exists := c.Get("userID")
	if exists {
		uid, _ := uuid.Parse(userIDStr.(string))
		_ = h.logoutUC.Execute(c.Request.Context(), uid)
	}

	// Clear cookies
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "code": "UNAUTHORIZED"})
		return
	}

	uid, _ := uuid.Parse(userIDStr.(string))
	user, err := h.meUC.Execute(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "code": "NOT_FOUND"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// setTokenCookies safely sets cookies with httpOnly
func (h *AuthHandler) setTokenCookies(c *gin.Context, tp *domain.TokenPair) {
	secure := false // change to true in prod against HTTPS
	if gin.Mode() == gin.ReleaseMode {
		secure = true
	}
	
	c.SetCookie("access_token", tp.AccessToken, int(15*time.Minute.Seconds()), "/", "", secure, true)
	c.SetCookie("refresh_token", tp.RefreshToken, int(7*24*time.Hour.Seconds()), "/", "", secure, true)
}
