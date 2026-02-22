package postgres

import (
	"context"
	"time"

	"github.com/auremsinistram/go-errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Connect(connString string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return errors.Wrap(err, "Postgres - Connect - #1")
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return errors.Wrap(err, "Postgres - Connect - #2")
	}

	p.Pool = pool

	if err := pool.Ping(ctx); err != nil {
		return errors.Wrap(err, "Postgres - Connect - #3")
	}

	return nil
}

func (p *Postgres) Close() error {
	p.Pool.Close()

	return nil
}
