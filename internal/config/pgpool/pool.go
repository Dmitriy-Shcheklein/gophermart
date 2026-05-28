// Package pgpool содержит логику инициализации пула для postgresql
package pgpool

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool структура для пула коннектов к БД
type Pool struct {
	// Pool пул коннектов
	Pool *pgxpool.Pool
}

// NewPool конструктор для пула
func NewPool(connString string) (*Pool, error) {
	pool := &Pool{}
	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return pool, err
	}

	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(
			&pgtype.Type{
				Name:  "uuid",
				OID:   pgtype.UUIDOID,
				Codec: &pgtype.UUIDCodec{},
			},
		)
		conn.TypeMap().RegisterDefaultPgType(&pgtype.UUID{}, "uuid")
		return nil
	}

	newPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return pool, err
	}
	pool.Pool = newPool

	return pool, nil
}
