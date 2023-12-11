package model

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Select(ctx context.Context, db Querier, dst interface{}, query string, args ...interface{}) error
}

type PgWorker struct {
	tx   *pgxpool.Tx
	pool *pgxpool.Pool
}

func NewPgWorker() (*PgWorker, error) {
	config, err := pgxpool.ParseConfig(config.Options.DatabaseURI)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 50

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}

	return &PgWorker{
		pool: pool}, nil
}

func (w *PgWorker) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	if w.tx != nil {
		return w.tx.Exec(ctx, sql, arguments...)
	}
	return w.pool.Exec(ctx, sql, arguments...)
}

func (w *PgWorker) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if w.tx != nil {
		return w.tx.Query(ctx, sql, args...)
	}
	return w.pool.Query(ctx, sql, args...)
}

func (w *PgWorker) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if w.tx != nil {
		return w.tx.QueryRow(ctx, sql, args...)
	}
	return w.pool.QueryRow(ctx, sql, args...)
}
func (w *PgWorker) Select(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	if w.tx != nil {
		return pgxscan.Select(ctx, w.tx, dst, query, args...)
	}
	return pgxscan.Select(ctx, w.pool, dst, query, args...)
}

func (w *PgWorker) Begin(ctx context.Context) error {
	if w.tx == nil {

		var tx *pgxpool.Tx

		t, err := w.pool.Begin(ctx)
		tx = t.(*pgxpool.Tx)

		if err != nil {
			return err
		}

		w.tx = tx
	}
	return nil
}

func (w *PgWorker) Rollback(ctx context.Context) error {
	if w.tx != nil {
		err := w.tx.Rollback(ctx)
		if err != nil {
			w.tx = nil
			return err
		}
	}
	return nil
}

func (w *PgWorker) Commit(ctx context.Context) error {
	if w.tx != nil {
		err := w.tx.Commit(ctx)
		if err != nil {
			w.tx = nil
			return err
		}
	}
	return nil
}
