package db

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"runtime"
	"strconv"
	"strings"

	"effectivemobile/pkg/gopkg-db/migrator"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(host string, port int, user, password, name string, options ...Option) (*DBPool, error) {
	pool := DBPool{
		size: runtime.NumCPU()*2 + 2,
	}
	for _, o := range options {
		o(&pool)
	}

	dsnList := map[string]string{
		"host":           host,
		"port":           strconv.Itoa(port),
		"user":           user,
		"password":       password,
		"dbname":         name,
		"sslmode":        "disable",
		"pool_max_conns": strconv.Itoa(pool.size),
	}
	dsnStr := make([]string, len(dsnList))
	for k, v := range dsnList {
		if v != "" && v != "0" {
			dsnStr = append(dsnStr, k+"="+v)
		}
	}
	dsn := strings.Join(dsnStr, " ")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}

	//dbLog := &dbLog{}
	//config.ConnConfig.Logger = dbLog

	pool.Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %w", err)
	}
	//dbLog.db = pool.Pool
	log.Printf("Use database %s@%s:%d/%s (conn size %v)\n", user, host, port, name, pool.size)

	// migrations
	mig, err := migrator.New(host, port, user, password, name)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to database: %w\n", err)
	}
	if err := mig.Up(); err != nil {
		return nil, fmt.Errorf("Can't roll up migrations: %w\n", err)
	}

	return &pool, nil
}

func ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	conn := FromContext(ctx)
	beginTx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	if err := fn(AddToContext(ctx, DBTx{beginTx})); err != nil {
		if err := beginTx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return beginTx.Commit(ctx)
}

func WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if f := flag.Lookup("test.v"); f != nil {
		return fn(ctx)
	}
	conn := FromContext(ctx)
	return conn.WithTx(ctx, fn)
}

type dbLog struct {
	db *pgxpool.Pool
}

func (l *dbLog) Log(ctx context.Context, level pgx.TxIsoLevel, msg string, data map[string]interface{}) {
	var ac string
	if l.db != nil {
		ac = fmt.Sprintf("%v/%v/%v", l.db.Stat().AcquiredConns(), l.db.Stat().TotalConns(), l.db.Stat().MaxConns())
	}
	log.Printf("[%v] %v %v %v\n", ac, level, data, msg)
}
