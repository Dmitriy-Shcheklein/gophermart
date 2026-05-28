package models

import "time"

// DbUser структура описывающая модель данных пользователя в БД
type DbUser struct {
	// ID идентификатор пользователя
	ID int `db:"id"`
	// Login логи пользователя
	Login string `db:"login"`
	// Password хэшированный пароль пользователя
	Password string `db:"password"`
	// CreatedAt время регистрация пользователя в системе
	CreatedAt time.Time `db:"created_at"`
}

// RegisterRequest структура описывающая модель данных ответа при регистрации юзера
type RegisterRequest struct {
	// Login логин пользователя
	Login string `json:"login" validate:"required,min=1"`
	// Password пароль
	Password string `json:"password" validate:"required,min=1"`
}

// AuthRequest структура описывающая модель данных ответа при авторизации юзера
type AuthRequest = RegisterRequest
