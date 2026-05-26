package pgpool

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	Pool *pgxpool.Pool
}

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
