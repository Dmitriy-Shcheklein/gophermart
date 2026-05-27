package models

type DbBalance struct {
	ID        int     `db:"id"`
	Current   float64 `db:"current"`
	Withdrawn float64 `db:"withdrawn"`
	UserID    int     `db:"user_id"`
}
