package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kurniawa9157/miniproject-interview/backend/internal/handler"
	"github.com/kurniawa9157/miniproject-interview/backend/internal/middleware"
	"github.com/kurniawa9157/miniproject-interview/backend/internal/repository"
	"github.com/kurniawa9157/miniproject-interview/backend/internal/service"
	"github.com/kurniawa9157/miniproject-interview/backend/pkg/database"
	googleoauth "github.com/kurniawa9157/miniproject-interview/backend/pkg/oauth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	ctx := context.Background()

	// Database
	db, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected")

	// Dependencies
	oauthConfig := googleoauth.NewGoogleConfig()
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, oauthConfig)
	authHandler := handler.NewAuthHandler(authService, userRepo, oauthConfig)

	// Router
	r := gin.Default()
	r.Use(middleware.CORS())

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.GET("/google", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", middleware.AuthRequired(), authHandler.Me)
	}

	// API routes (protected)
	api := r.Group("/api", middleware.AuthRequired())
	{
		// Orders — Phase 2
		_ = api
	}

	// Admin routes — Phase 3
	admin := r.Group("/api/admin", middleware.AuthRequired(), middleware.AdminRequired())
	{
		_ = admin
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
