package models

import (
	"database/sql"
	"time"
)

type (
	// OrderStatuses тип описывающий статусы заказов
	OrderStatuses string
	// LoyaltyStatuses тип описывающий статусы расчетов
	LoyaltyStatuses string
)

const (
	// NewOrder заказ загружен в систему, но не попал в обработку
	NewOrder OrderStatuses = "NEW"
	// ProcessingOrder вознаграждение за заказ рассчитывается
	ProcessingOrder OrderStatuses = "PROCESSING"
	// InvalidOrder система расчёта вознаграждений отказала в расчёте
	InvalidOrder OrderStatuses = "INVALID"
	// ProcessedOrder данные по заказу проверены и информация о расчёте успешно получена
	ProcessedOrder OrderStatuses = "PROCESSED"
)

const (
	// Registered заказ зарегистрирован, но вознаграждение не рассчитано
	Registered = "REGISTERED"
	// Invalid заказ не принят к расчёту, и вознаграждение не будет начислено
	Invalid = "INVALID"
	// Processing расчёт начисления в процессе
	Processing = "PROCESSING"
	// Processed расчёт начисления окончен
	Processed = "PROCESSED"
)

// DbOrder структура описывающая модель данных заказа в БД
type DbOrder struct {
	// ID идентификатора записи в бд
	ID int `db:"id"`
	// Number номер заказа
	Number string `db:"number"`
	// Status статус заказа
	Status OrderStatuses `db:"status"`
	// UploadedAt время загрузки заказа
	UploadedAt time.Time `db:"uploaded_at"`
	// Accrual количество баллов за заказ
	Accrual sql.NullFloat64 `db:"accrual"`
	// UserID идентификатор пользователя
	UserID int `db:"user_id"`
}

// RequestOrder структура описывающая модель данных ответа по заказу
type RequestOrder struct {
	// Number номер заказа
	Number string `json:"number"`
	// Status статус заказа
	Status OrderStatuses `json:"status"`
	// UploadedAt время загрузки заказа
	UploadedAt time.Time `json:"uploaded_at"`
	// Accrual количество баллов за заказ
	Accrual *float64 `json:"accrual,omitempty"`
}

// LoyaltyOrderData структура описывающая модель данных ответа от системы расчета
type LoyaltyOrderData struct {
	// Order идентификатор заказа
	Order string `json:"order" validate:"required,min=1"`
	// Status статус расчета
	Status LoyaltyStatuses `json:"status" validate:"required,oneof=REGISTERED INVALID PROCESSING PROCESSED"`
	// Accrual количество баллов за заказ
	Accrual *float64 `json:"accrual,omitempty"`
}
