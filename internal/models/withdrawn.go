package models

import "time"

// DbWithdrawn структура описывающая модель данных по списанию в БД
type DbWithdrawn struct {
	// ID идентификатор записи в БД
	ID int `db:"id"`
	// Sum количество баллов
	Sum float64 `db:"sum"`
	// Order номер заказа
	Order string `db:"order_num"`
	// UserID идентификатор пользователя
	UserID int `db:"user_id"`
	// ProcessedAt время начала списания
	ProcessedAt time.Time `db:"processed_at"`
}

// RequestWithdrawn структура описывающая запрос на списание
type RequestWithdrawn struct {
	// Sum количество баллов
	Sum float64 `json:"sum" validate:"required,min=0.01"`
	// Order номер заказа
	Order string `json:"order" validate:"required,min=1"`
}

// ResponseWithdrawn структура описывающая ответ на списание
type ResponseWithdrawn struct {
	// Sum количество баллов
	Sum float64 `json:"sum"`
	// Order номер заказа
	Order string `json:"order"`
	// ProcessedAt время начала списания
	ProcessedAt time.Time `json:"processed_at"`
}
