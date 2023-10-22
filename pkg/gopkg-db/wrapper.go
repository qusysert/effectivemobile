package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IClient interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error

	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// DB POOL

type DBPool struct {
	*pgxpool.Pool
	size int
}

func (d *DBPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return d.Pool.Query(ctx, sql, args...)
}

func (d *DBPool) Commit(ctx context.Context) error {
	return nil
}
func (d *DBPool) Rollback(ctx context.Context) error {
	return nil
}

func (d *DBPool) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	ctxWithTx := AddToContext(ctx, DBTx{tx})
	if err := fn(ctxWithTx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

// DB TRANSACTION

type DBTx struct {
	pgx.Tx
}

func (tx DBTx) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return tx.Tx, nil
}

func (tx DBTx) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
