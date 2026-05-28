// Package models пакет содержащий модели БД, запросов и ответов
package models

// DbBalance структура описывающая модель данных баланса в БД
type DbBalance struct {
	// ID идентификатор строки в БД
	ID int `db:"id"`
	// Current текущий баланс
	Current float64 `db:"current"`
	// Withdrawn всего баллов за все время
	Withdrawn float64 `db:"withdrawn"`
	// UserID идентификатор пользователя
	UserID int `db:"user_id"`
}

// ResponseBalance структура описывающая модель ответа данных по балансу
type ResponseBalance struct {
	// Current текущий баланс
	Current float64 `json:"current"`
	// Withdrawn всего баллов за все время
	Withdrawn float64 `json:"withdrawn"`
}
