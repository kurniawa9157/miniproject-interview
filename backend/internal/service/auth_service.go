package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	googleoauth "jumpapay/backend/pkg/oauth"
	"jumpapay/backend/internal/model"
	"jumpapay/backend/internal/repository"
	"golang.org/x/oauth2"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	oauthConfig *oauth2.Config
}

func NewAuthService(userRepo *repository.UserRepository, oauthConfig *oauth2.Config) *AuthService {
	return &AuthService{userRepo: userRepo, oauthConfig: oauthConfig}
}

// HandleCallback exchanges OAuth code for user info, upserts user, returns JWT.
func (s *AuthService) HandleCallback(ctx context.Context, code string) (string, *model.User, error) {
	info, err := googleoauth.GetUserInfo(ctx, s.oauthConfig, code)
	if err != nil {
		return "", nil, fmt.Errorf("get google user info: %w", err)
	}

	user, err := s.userRepo.UpsertByGoogleID(ctx, info.ID, info.Name, info.Email, info.Picture)
	if err != nil {
		return "", nil, fmt.Errorf("upsert user: %w", err)
	}

	// Override is_admin from env whitelist
	user.IsAdmin = isAdminEmail(user.Email)

	token, err := generateJWT(user)
	if err != nil {
		return "", nil, fmt.Errorf("generate jwt: %w", err)
	}

	return token, user, nil
}

type JWTClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func generateJWT(user *model.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	claims := JWTClaims{
		UserID:  user.ID,
		Email:   user.Email,
		IsAdmin: isAdminEmail(user.Email),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenStr string) (*JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func isAdminEmail(email string) bool {
	adminEmails := os.Getenv("ADMIN_EMAILS")
	if adminEmails == "" {
		return false
	}
	for _, e := range strings.Split(adminEmails, ",") {
		if strings.TrimSpace(e) == email {
			return true
		}
	}
	return false
}
