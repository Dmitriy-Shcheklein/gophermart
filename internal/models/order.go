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
	Status     string          `db:"status"`
	UploadedAt time.Time       `db:"uploaded_at"`
	Accrual    sql.NullFloat64 `db:"accrual"`
	UserID     int             `db:"user_id"`
}
