package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"jumpapay/backend/internal/middleware"
	"jumpapay/backend/internal/repository"
	"jumpapay/backend/internal/service"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	authService *service.AuthService
	userRepo    *repository.UserRepository
	oauthConfig *oauth2.Config
}

func NewAuthHandler(authService *service.AuthService, userRepo *repository.UserRepository, oauthConfig *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
		oauthConfig: oauthConfig,
	}
}

// GET /auth/google — redirect ke Google consent screen
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /auth/google/callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	token, user, err := h.authService.HandleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	secure := os.Getenv("APP_ENV") == "production"
	c.SetCookie("token", token, int((7 * 24 * time.Hour).Seconds()), "/", "", secure, true)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	// Redirect ke frontend, bawa info minimal via query jika perlu
	if user.IsAdmin {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/admin")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/orders/new")
	}
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// GET /auth/me — requires AuthRequired middleware
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	user, err := h.userRepo.FindByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"photo_url":  user.PhotoURL,
		"is_admin":   service.IsAdminEmail(user.Email),
		"created_at": user.CreatedAt,
	})
}
