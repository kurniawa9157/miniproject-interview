package service

import (
	"context"
	"errors"
	"fmt"

	"jumpapay/backend/internal/model"
	"jumpapay/backend/internal/repository"
)

var ErrOrderNotFound = errors.New("order not found")
var ErrInvalidTransition = errors.New("invalid status transition")

type AdminService struct {
	orderRepo *repository.OrderRepository
}

func NewAdminService(orderRepo *repository.OrderRepository) *AdminService {
	return &AdminService{orderRepo: orderRepo}
}

func (s *AdminService) ListOrders(ctx context.Context, statusFilter string) ([]model.OrderWithLogs, error) {
	return s.orderRepo.ListAll(ctx, statusFilter)
}

func (s *AdminService) GetOrderDetail(ctx context.Context, orderID string) (*model.OrderWithLogs, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID, "")
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *AdminService) UpdateStatus(ctx context.Context, orderID, adminUserID string, newStatus model.OrderStatus) error {
	order, err := s.orderRepo.FindByID(ctx, orderID, "")
	if err != nil {
		return ErrOrderNotFound
	}

	if !order.Status.CanTransitionTo(newStatus) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, order.Status, newStatus)
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, adminUserID, newStatus)
}
