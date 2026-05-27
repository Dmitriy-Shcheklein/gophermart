package models

import (
	"database/sql"
	"time"
)

type (
	OrderStatuses   string
	LoyaltyStatuses string
)

const (
	NewOrder        OrderStatuses = "NEW"
	ProcessingOrder OrderStatuses = "PROCESSING"
	InvalidOrder    OrderStatuses = "INVALID"
	ProcessedOrder  OrderStatuses = "PROCESSED"
)

const (
	Registered = "REGISTERED"
	Invalid    = "INVALID"
	Processing = "PROCESSING"
	Processed  = "PROCESSED"
)

type DbOrder struct {
	ID         int             `db:"id"`
	Number     string          `db:"number"`
	Status     OrderStatuses   `db:"status"`
	UploadedAt time.Time       `db:"uploaded_at"`
	Accrual    sql.NullFloat64 `db:"accrual"`
	UserID     int             `db:"user_id"`
}

type RequestOrder struct {
	Number     string        `json:"number"`
	Status     OrderStatuses `json:"status"`
	UploadedAt time.Time     `json:"uploaded_at"`
	Accrual    *float64      `json:"accrual,omitempty"`
}

type LoyaltyOrderData struct {
	Order   string          `json:"order" validate:"required,min=1"`
	Status  LoyaltyStatuses `json:"status" validate:"required,oneof=REGISTERED INVALID PROCESSING PROCESSED"`
	Accrual *float64        `json:"accrual,omitempty"`
}
