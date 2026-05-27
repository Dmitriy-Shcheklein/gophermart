package models

type DbWithdrawn struct {
	ID     int     `db:"id"`
	Sum    float64 `db:"sum"`
	Order  string  `db:"order_num"`
	UserID int     `db:"user_id"`
}

type RequestWithdraw struct {
	Sum   float64 `json:"sum" validate:"required,min=0.01"`
	Order string  `json:"order" validate:"required,min=1"`
}
