package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"jumpapay/backend/internal/middleware"
	"jumpapay/backend/internal/service"
	"jumpapay/backend/pkg/payment"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// GET /api/payment/config — return Midtrans client key for Snap.js
func (h *PaymentHandler) Config(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"client_key": h.paymentService.ClientKey()})
}

// POST /api/orders/:id/pay — create Snap token (customer, own order only)
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString(middleware.UserIDKey)

	token, err := h.paymentService.CreateSnapToken(c.Request.Context(), orderID, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "order tidak ditemukan"})
		case errors.Is(err, service.ErrAlreadyPaid):
			c.JSON(http.StatusConflict, gin.H{"error": "order sudah dibayar"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat transaksi pembayaran"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// POST /payment/notification — Midtrans webhook (public, verified by signature)
func (h *PaymentHandler) Notification(c *gin.Context) {
	var n payment.Notification
	if err := c.ShouldBindJSON(&n); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if err := h.paymentService.HandleNotification(c.Request.Context(), n); err != nil {
		if errors.Is(err, service.ErrInvalidSignature) {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
			return
		}
		log.Printf("payment notification failed (order=%s status=%s): %v", n.OrderID, n.TransactionStatus, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memproses notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
