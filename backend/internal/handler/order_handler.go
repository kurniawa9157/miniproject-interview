package handler

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"jumpapay/backend/internal/middleware"
	"jumpapay/backend/internal/model"
	"jumpapay/backend/internal/service"
)

var (
	reWA    = regexp.MustCompile(`^(\+62|08)\d{8,13}$`)
	rePlate = regexp.MustCompile(`^[A-Z]{1,2}\s\d{1,4}\s[A-Z]{1,3}$`)
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// POST /api/orders
func (h *OrderHandler) Submit(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)

	whatsapp := strings.TrimSpace(c.PostForm("whatsapp"))
	plateNumber := strings.ToUpper(strings.TrimSpace(c.PostForm("plate_number")))
	frameNumber := strings.ToUpper(strings.TrimSpace(c.PostForm("frame_number")))

	// Validate
	if whatsapp == "" || plateNumber == "" || frameNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "whatsapp, plate_number, dan frame_number wajib diisi"})
		return
	}
	if !reWA.MatchString(whatsapp) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format nomor WhatsApp tidak valid (contoh: 08123456789 atau +6281234567890)"})
		return
	}
	if !rePlate.MatchString(plateNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format nomor plat tidak valid (contoh: D 1234 ABC)"})
		return
	}
	if len(frameNumber) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nomor rangka harus tepat 5 karakter"})
		return
	}

	ktpFile, ktpHeader, err := c.Request.FormFile("ktp")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file KTP wajib diupload"})
		return
	}
	defer ktpFile.Close()

	stnkFile, stnkHeader, err := c.Request.FormFile("stnk")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file STNK wajib diupload"})
		return
	}
	defer stnkFile.Close()

	// File size check (2MB)
	const maxSize = 2 << 20
	if ktpHeader.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ukuran file KTP melebihi 2MB"})
		return
	}
	if stnkHeader.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ukuran file STNK melebihi 2MB"})
		return
	}

	// File type check
	if !isImageFile(ktpHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file KTP harus berformat JPG atau PNG"})
		return
	}
	if !isImageFile(stnkHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file STNK harus berformat JPG atau PNG"})
		return
	}

	order, err := h.orderService.Submit(c.Request.Context(), service.SubmitOrderInput{
		UserID:      userID,
		Whatsapp:    whatsapp,
		PlateNumber: plateNumber,
		FrameNumber: frameNumber,
		KtpFile:     ktpFile,
		KtpHeader:   ktpHeader,
		StnkFile:    stnkFile,
		StnkHeader:  stnkHeader,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": order})
}

// GET /api/orders
func (h *OrderHandler) ListMine(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)

	orders, err := h.orderService.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data order"})
		return
	}

	if orders == nil {
		orders = []model.Order{}
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GET /api/orders/:id
func (h *OrderHandler) GetTracking(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey)
	orderID := c.Param("id")

	order, err := h.orderService.GetByID(c.Request.Context(), orderID, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "order tidak ditemukan"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "order tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func isImageFile(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".jpg") ||
		strings.HasSuffix(lower, ".jpeg") ||
		strings.HasSuffix(lower, ".png")
}
