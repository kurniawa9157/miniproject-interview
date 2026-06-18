package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"jumpapay/backend/internal/service"
)

const UserIDKey = "user_id"
const UserEmailKey = "user_email"
const IsAdminKey = "is_admin"

// AuthRequired validates JWT from cookie or Authorization header.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := tokenFromRequest(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, err := service.ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(IsAdminKey, claims.IsAdmin)
		c.Next()
	}
}

// AdminRequired must be chained after AuthRequired.
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, _ := c.Get(IsAdminKey)
		if isAdmin != true {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func tokenFromRequest(c *gin.Context) string {
	// 1. httpOnly cookie
	if cookie, err := c.Cookie("token"); err == nil && cookie != "" {
		return cookie
	}
	// 2. Authorization: Bearer <token>
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
