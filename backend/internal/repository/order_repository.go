package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"jumpapay/backend/internal/model"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// GenerateID creates order ID in format JP-YYYYMMDD-XXXX.
func (r *OrderRepository) GenerateID(ctx context.Context) (string, error) {
	today := time.Now().Format("20060102")
	prefix := fmt.Sprintf("JP-%s-", today)

	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM orders WHERE id LIKE $1`,
		prefix+"%",
	).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("count orders today: %w", err)
	}

	return fmt.Sprintf("%s%04d", prefix, count+1), nil
}

// Create inserts a new order and its initial PENDING status log.
func (r *OrderRepository) Create(ctx context.Context, o *model.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (id, user_id, whatsapp, plate_number, frame_number, ktp_url, stnk_url, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, o.ID, o.UserID, o.Whatsapp, o.PlateNumber, o.FrameNumber, o.KtpURL, o.StnkURL, o.Status)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO order_status_logs (order_id, status, changed_by)
		VALUES ($1, $2, $3)
	`, o.ID, o.Status, o.UserID)
	if err != nil {
		return fmt.Errorf("insert status log: %w", err)
	}

	return tx.Commit(ctx)
}

// FindByID returns order + status logs. Optionally enforce ownerUserID (empty = skip check).
func (r *OrderRepository) FindByID(ctx context.Context, id, ownerUserID string) (*model.OrderWithLogs, error) {
	row := r.db.QueryRow(ctx, `
		SELECT o.id, o.user_id, o.whatsapp, o.plate_number, o.frame_number,
		       o.ktp_url, o.stnk_url, o.status, o.created_at, o.updated_at,
		       u.name, u.email
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE o.id = $1
	`, id)

	var o model.OrderWithLogs
	err := row.Scan(
		&o.ID, &o.UserID, &o.Whatsapp, &o.PlateNumber, &o.FrameNumber,
		&o.KtpURL, &o.StnkURL, &o.Status, &o.CreatedAt, &o.UpdatedAt,
		&o.UserName, &o.UserEmail,
	)
	if err != nil {
		return nil, fmt.Errorf("find order: %w", err)
	}

	if ownerUserID != "" && o.UserID != ownerUserID {
		return nil, fmt.Errorf("forbidden")
	}

	logs, err := r.findLogs(ctx, id)
	if err != nil {
		return nil, err
	}
	o.StatusLogs = logs

	return &o, nil
}

// FindByUserID returns all orders belonging to a user.
func (r *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]model.Order, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, whatsapp, plate_number, frame_number,
		       ktp_url, stnk_url, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("find orders by user: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.Whatsapp, &o.PlateNumber, &o.FrameNumber,
			&o.KtpURL, &o.StnkURL, &o.Status, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// ListAll returns all orders with optional status filter (empty = all).
func (r *OrderRepository) ListAll(ctx context.Context, statusFilter string) ([]model.OrderWithLogs, error) {
	query := `
		SELECT o.id, o.user_id, o.whatsapp, o.plate_number, o.frame_number,
		       o.ktp_url, o.stnk_url, o.status, o.created_at, o.updated_at,
		       u.name, u.email
		FROM orders o
		JOIN users u ON u.id = o.user_id
	`
	args := []any{}
	if statusFilter != "" {
		query += " WHERE o.status = $1"
		args = append(args, statusFilter)
	}
	query += " ORDER BY o.created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list all orders: %w", err)
	}
	defer rows.Close()

	var orders []model.OrderWithLogs
	for rows.Next() {
		var o model.OrderWithLogs
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.Whatsapp, &o.PlateNumber, &o.FrameNumber,
			&o.KtpURL, &o.StnkURL, &o.Status, &o.CreatedAt, &o.UpdatedAt,
			&o.UserName, &o.UserEmail,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// UpdateStatus changes order status and inserts a log entry.
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID, changedByUserID string, newStatus model.OrderStatus) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2
	`, newStatus, orderID)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO order_status_logs (order_id, status, changed_by)
		VALUES ($1, $2, $3)
	`, orderID, newStatus, changedByUserID)
	if err != nil {
		return fmt.Errorf("insert status log: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) findLogs(ctx context.Context, orderID string) ([]model.OrderStatusLog, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, order_id, status, changed_by, created_at
		FROM order_status_logs
		WHERE order_id = $1
		ORDER BY created_at ASC
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("find logs: %w", err)
	}
	defer rows.Close()

	var logs []model.OrderStatusLog
	for rows.Next() {
		var l model.OrderStatusLog
		if err := rows.Scan(&l.ID, &l.OrderID, &l.Status, &l.ChangedBy, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
