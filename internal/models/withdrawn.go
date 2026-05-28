package models

import "time"

type DbWithdrawn struct {
	ID          int       `db:"id"`
	Sum         float64   `db:"sum"`
	Order       string    `db:"order_num"`
	UserID      int       `db:"user_id"`
	ProcessedAt time.Time `db:"processed_at"`
}

type RequestWithdrawn struct {
	Sum   float64 `json:"sum" validate:"required,min=0.01"`
	Order string  `json:"order" validate:"required,min=1"`
}

type ResponseWithdrawn struct {
	Sum         float64   `json:"sum"`
	Order       string    `json:"order"`
	ProcessedAt time.Time `json:"processed_at"`
}
