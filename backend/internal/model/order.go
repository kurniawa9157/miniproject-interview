package model

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusInProcess OrderStatus = "IN_PROCESS"
	StatusDone      OrderStatus = "DONE"
	StatusCancelled OrderStatus = "CANCELLED"
)

// ValidTransitions maps current status to allowed next statuses.
var ValidTransitions = map[OrderStatus][]OrderStatus{
	StatusPending:   {StatusInProcess, StatusCancelled},
	StatusInProcess: {StatusDone, StatusCancelled},
	StatusDone:      {},
	StatusCancelled: {},
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	for _, allowed := range ValidTransitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

type PaymentStatus string

const (
	PaymentUnpaid  PaymentStatus = "UNPAID"
	PaymentPending PaymentStatus = "PENDING"
	PaymentPaid    PaymentStatus = "PAID"
	PaymentFailed  PaymentStatus = "FAILED"
)

type Order struct {
	ID            string        `json:"id"`
	UserID        string        `json:"user_id"`
	Whatsapp      string        `json:"whatsapp"`
	PlateNumber   string        `json:"plate_number"`
	FrameNumber   string        `json:"frame_number"`
	KtpURL        string        `json:"ktp_url"`
	StnkURL       string        `json:"stnk_url"`
	Status        OrderStatus   `json:"status"`
	Amount        int           `json:"amount"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type OrderStatusLog struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Status    OrderStatus `json:"status"`
	ChangedBy *string   `json:"changed_by"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderWithLogs struct {
	Order
	StatusLogs []OrderStatusLog `json:"status_logs"`
	UserName   string           `json:"user_name"`
	UserEmail  string           `json:"user_email"`
}
