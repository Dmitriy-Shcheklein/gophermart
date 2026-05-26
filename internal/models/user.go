package models

import "time"

type DbUser struct {
	ID        int       `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required,min=1"`
	Password string `json:"password" validate:"required,min=1"`
}

type AuthRequest = RegisterRequest
