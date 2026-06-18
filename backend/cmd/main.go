package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"jumpapay/backend/internal/handler"
	"jumpapay/backend/internal/middleware"
	"jumpapay/backend/internal/repository"
	"jumpapay/backend/internal/service"
	"jumpapay/backend/pkg/database"
	googleoauth "jumpapay/backend/pkg/oauth"
	"jumpapay/backend/pkg/payment"
	"jumpapay/backend/pkg/storage"
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

	// Storage
	storageClient, err := storage.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}
	log.Println("Storage connected")

	// Dependencies
	oauthConfig := googleoauth.NewGoogleConfig()
	midtransClient := payment.NewClient()

	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	authService := service.NewAuthService(userRepo, oauthConfig)
	orderService := service.NewOrderService(orderRepo, storageClient)
	adminService := service.NewAdminService(orderRepo)
	paymentService := service.NewPaymentService(orderRepo, midtransClient)

	authHandler := handler.NewAuthHandler(authService, userRepo, oauthConfig)
	orderHandler := handler.NewOrderHandler(orderService)
	adminHandler := handler.NewAdminHandler(adminService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Router
	r := gin.Default()
	r.Use(middleware.CORS())

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Midtrans webhook (public - dipanggil oleh Midtrans, diverifikasi via signature)
	r.POST("/payment/notification", paymentHandler.Notification)

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.GET("/google", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", middleware.AuthRequired(), authHandler.Me)
	}

	// Customer API routes
	api := r.Group("/api", middleware.AuthRequired())
	{
		api.POST("/orders", orderHandler.Submit)
		api.GET("/orders", orderHandler.ListMine)
		api.GET("/orders/:id", orderHandler.GetTracking)

		// Payment (bonus)
		api.GET("/payment/config", paymentHandler.Config)
		api.POST("/orders/:id/pay", paymentHandler.CreatePayment)
	}

	// Admin routes
	admin := r.Group("/api/admin", middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.GET("/orders", adminHandler.ListOrders)
		admin.GET("/orders/:id", adminHandler.GetOrderDetail)
		admin.PATCH("/orders/:id/status", adminHandler.UpdateStatus)
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