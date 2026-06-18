package service

import (
	"context"
	"errors"
	"fmt"

	"jumpapay/backend/internal/model"
	"jumpapay/backend/internal/repository"
	"jumpapay/backend/pkg/payment"
)

var ErrInvalidSignature = errors.New("invalid signature")
var ErrAlreadyPaid = errors.New("order already paid")

type PaymentService struct {
	orderRepo *repository.OrderRepository
	midtrans  *payment.Client
}

func NewPaymentService(orderRepo *repository.OrderRepository, midtrans *payment.Client) *PaymentService {
	return &PaymentService{orderRepo: orderRepo, midtrans: midtrans}
}

func (s *PaymentService) ClientKey() string {
	return s.midtrans.ClientKey()
}

// CreateSnapToken creates a Midtrans transaction for an order owned by ownerUserID
// and returns the Snap token to render the payment popup.
func (s *PaymentService) CreateSnapToken(ctx context.Context, orderID, ownerUserID string) (string, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID, ownerUserID)
	if err != nil {
		return "", ErrOrderNotFound
	}

	if order.PaymentStatus == model.PaymentPaid {
		return "", ErrAlreadyPaid
	}

	token, err := s.midtrans.CreateTransaction(payment.SnapRequest{
		OrderID:      order.ID,
		Amount:       order.Amount,
		CustomerName: order.UserName,
		Email:        order.UserEmail,
	})
	if err != nil {
		return "", fmt.Errorf("create transaction: %w", err)
	}

	if err := s.orderRepo.SetPaymentToken(ctx, order.ID, token); err != nil {
		return "", err
	}

	return token, nil
}

// HandleNotification processes a Midtrans webhook notification.
func (s *PaymentService) HandleNotification(ctx context.Context, n payment.Notification) error {
	if !s.midtrans.VerifySignature(n) {
		return ErrInvalidSignature
	}

	switch {
	case n.IsSuccess():
		return s.orderRepo.MarkPaid(ctx, n.OrderID)
	case n.IsFailure():
		return s.orderRepo.MarkPaymentFailed(ctx, n.OrderID)
	default:
		// pending / other — no state change
		return nil
	}
}
