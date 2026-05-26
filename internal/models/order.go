package models

import (
	"database/sql"
	"time"
)

type OrderStatuses string

const (
	NewOrder        OrderStatuses = "NEW"
	ProcessingOrder OrderStatuses = "PROCESSING"
	InvalidOrder    OrderStatuses = "INVALID"
	ProcessedOrder  OrderStatuses = "PROCESSED"
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
