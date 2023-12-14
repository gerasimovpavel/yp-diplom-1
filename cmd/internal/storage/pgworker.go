package storage

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContextKey string

type PgWorker struct {
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
	if ctx.Value(ContextKey("tx")) != nil {
		tx, ok := ctx.Value(ContextKey("tx")).(*pgxpool.Tx)
		if !ok {
			panic(errors.New("not ok"))
		}
		return tx.Exec(ctx, sql, arguments...)
	}
	return w.pool.Exec(ctx, sql, arguments...)
}

func (w *PgWorker) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if ctx.Value(ContextKey("tx")) != nil {
		return ctx.Value(ContextKey("tx")).(*pgxpool.Tx).Query(ctx, sql, args...)
	}
	return w.pool.Query(ctx, sql, args...)
}

func (w *PgWorker) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if ctx.Value(ContextKey("tx")) != nil {
		return ctx.Value(ContextKey("tx")).(*pgxpool.Tx).QueryRow(ctx, sql, args...)
	}
	return w.pool.QueryRow(ctx, sql, args...)
}
func (w *PgWorker) Select(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	if ctx.Value(ContextKey("tx")) != nil {
		tx, ok := ctx.Value(ContextKey("tx")).(*pgxpool.Tx)
		if !ok {
			panic(errors.New("not ok"))
		}
		return pgxscan.Select(ctx, tx, dst, query, args...)
	}
	return pgxscan.Select(ctx, w.pool, dst, query, args...)
}

func (w *PgWorker) Begin(ctx context.Context) (context.Context, error) {
	if ctx.Value(ContextKey("tx")) == nil {
		t, err := w.pool.Begin(ctx)
		if err != nil {
			logger.Logger.Sugar().Errorln(err)
			return ctx, err
		}
		ctx = context.WithValue(ctx, ContextKey("tx"), t.(*pgxpool.Tx))
		return ctx, nil
	}
	return ctx, nil
}

func (w *PgWorker) Rollback(ctx context.Context) error {
	t := ctx.Value(ContextKey("tx")).(*pgxpool.Tx)
	if t != nil {
		err := t.Rollback(ctx)
		if err != nil {
			logger.Logger.Sugar().Errorln(err)
			return err
		}
	}
	return nil
}

func (w *PgWorker) Commit(ctx context.Context) error {
	t := ctx.Value(ContextKey("tx")).(*pgxpool.Tx)
	if t != nil {
		err := t.Commit(ctx)
		if err != nil {
			logger.Logger.Sugar().Errorln(err)
			return err
		}
	}
	return nil
}
