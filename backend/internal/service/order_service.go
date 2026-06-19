package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"jumpapay/backend/internal/model"
	"jumpapay/backend/internal/repository"
	"jumpapay/backend/pkg/storage"
)

// DefaultOrderAmount is the hardcoded STNK renewal fee (Rp 150.000) per the spec.
const DefaultOrderAmount = 150000

type SubmitOrderInput struct {
	UserID      string
	Whatsapp    string
	PlateNumber string
	FrameNumber string
	KtpFile     multipart.File
	KtpHeader   *multipart.FileHeader
	StnkFile    multipart.File
	StnkHeader  *multipart.FileHeader
}

type OrderService struct {
	orderRepo *repository.OrderRepository
	storage   *storage.Client
}

func NewOrderService(orderRepo *repository.OrderRepository, storage *storage.Client) *OrderService {
	return &OrderService{orderRepo: orderRepo, storage: storage}
}

func (s *OrderService) Submit(ctx context.Context, input SubmitOrderInput) (*model.Order, error) {
	orderID, err := s.orderRepo.GenerateID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate order id: %w", err)
	}

	ktpURL, err := s.uploadFile(ctx, orderID, "ktp", input.KtpFile, input.KtpHeader)
	if err != nil {
		return nil, fmt.Errorf("upload ktp: %w", err)
	}

	stnkURL, err := s.uploadFile(ctx, orderID, "stnk", input.StnkFile, input.StnkHeader)
	if err != nil {
		return nil, fmt.Errorf("upload stnk: %w", err)
	}

	order := &model.Order{
		ID:            orderID,
		UserID:        input.UserID,
		Whatsapp:      input.Whatsapp,
		PlateNumber:   input.PlateNumber,
		FrameNumber:   strings.ToUpper(input.FrameNumber),
		KtpURL:        ktpURL,
		StnkURL:       stnkURL,
		Status:        model.StatusPending,
		Amount:        DefaultOrderAmount,
		PaymentStatus: model.PaymentUnpaid,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetByID(ctx context.Context, orderID, ownerUserID string) (*model.OrderWithLogs, error) {
	return s.orderRepo.FindByID(ctx, orderID, ownerUserID)
}

func (s *OrderService) ListByUser(ctx context.Context, userID string) ([]model.Order, error) {
	return s.orderRepo.FindByUserID(ctx, userID)
}

func (s *OrderService) uploadFile(ctx context.Context, orderID, docType string, file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	objectName := fmt.Sprintf("orders/%s/%s%s", orderID, docType, ext)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// Seek to beginning in case it was read before
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	url, err := s.storage.UploadFile(ctx, objectName, contentType, file, header.Size)
	if err != nil {
		return "", err
	}

	return url, nil
}
