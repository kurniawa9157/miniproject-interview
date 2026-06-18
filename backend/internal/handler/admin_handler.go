package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"jumpapay/backend/internal/middleware"
	"jumpapay/backend/internal/model"
	"jumpapay/backend/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GET /api/admin/orders?status=PENDING
func (h *AdminHandler) ListOrders(c *gin.Context) {
	statusFilter := c.Query("status")

	if statusFilter != "" {
		s := model.OrderStatus(statusFilter)
		if s != model.StatusPending && s != model.StatusInProcess &&
			s != model.StatusDone && s != model.StatusCancelled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status tidak valid"})
			return
		}
	}

	orders, err := h.adminService.ListOrders(c.Request.Context(), statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data order"})
		return
	}

	if orders == nil {
		orders = []model.OrderWithLogs{}
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GET /api/admin/orders/:id
func (h *AdminHandler) GetOrderDetail(c *gin.Context) {
	orderID := c.Param("id")

	order, err := h.adminService.GetOrderDetail(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// PATCH /api/admin/orders/:id/status
func (h *AdminHandler) UpdateStatus(c *gin.Context) {
	orderID := c.Param("id")
	adminUserID := c.GetString(middleware.UserIDKey)

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field status wajib diisi"})
		return
	}

	newStatus := model.OrderStatus(body.Status)
	if newStatus != model.StatusInProcess && newStatus != model.StatusDone && newStatus != model.StatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nilai status tidak valid"})
		return
	}

	if err := h.adminService.UpdateStatus(c.Request.Context(), orderID, adminUserID, newStatus); err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order tidak ditemukan"})
			return
		}
		if errors.Is(err, service.ErrInvalidTransition) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengupdate status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status berhasil diupdate"})
}
