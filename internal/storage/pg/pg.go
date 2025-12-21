package pg

import (
	"context"

	"github.com/pkg/errors"

	"github.com/dmnAlex/gophermart/internal/model/errx"
	"github.com/golang-migrate/migrate/v4"
	pgxdriver "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	stopCtx context.Context
	pool    *pgxpool.Pool
	tx      pgx.Tx
}

func New(ctx context.Context, dsn, migrationsPath string) (*DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "create new pool")
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errors.Wrap(err, "ping pool")
	}
	if err := applyMigrations(pool, migrationsPath); err != nil {
		pool.Close()
		return nil, errors.Wrap(err, "apply migrations")
	}
	return &DB{stopCtx: ctx, pool: pool}, nil
}

func applyMigrations(pool *pgxpool.Pool, migrationsPath string) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	driver, err := pgxdriver.WithInstance(sqlDB, &pgxdriver.Config{})
	if err != nil {
		return errors.Wrap(err, "get driver with instance")
	}
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "migrate with db instance")
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "up migrations")
	}

	return nil
}

func (db *DB) Close() error {
	if db.tx != nil {
		db.tx.Rollback(db.stopCtx)
	}
	db.pool.Close()

	return nil
}

func (db *DB) WithCtx(ctx context.Context) *DB {
	return &DB{
		stopCtx: ctx,
		pool:    db.pool,
		tx:      db.tx,
	}
}

func (db *DB) DoTx(f func(*DB) error, opts ...*pgx.TxOptions) error {
	if db.tx != nil {
		return errors.Wrap(f(db), "call func with tx")
	}
	var opt *pgx.TxOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt == nil {
		opt = &pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		}
	}
	tx, err := db.pool.BeginTx(db.stopCtx, *opt)
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	txDB := &DB{
		stopCtx: db.stopCtx,
		pool:    db.pool,
		tx:      tx,
	}
	if err := f(txDB); err != nil {
		tx.Rollback(db.stopCtx)
		return errors.Wrap(err, "call func")
	}
	return errors.Wrap(tx.Commit(db.stopCtx), "commit")
}

func (db *DB) Exec(query string, args ...any) (pgconn.CommandTag, error) {
	if db.tx != nil {
		return db.tx.Exec(db.stopCtx, query, args...)
	}
	return db.pool.Exec(db.stopCtx, query, args...)
}

func (db *DB) Query(query string, args pgx.NamedArgs, f func(pgx.Rows) error) error {
	var (
		rows pgx.Rows
		err  error
	)
	if db.tx != nil {
		rows, err = db.tx.Query(db.stopCtx, query, args)
	} else {
		rows, err = db.pool.Query(db.stopCtx, query, args)
	}
	if err != nil {
		return errors.Wrap(err, "query")
	}
	defer rows.Close()
	for rows.Next() {
		if err := f(rows); err != nil {
			return errors.Wrap(err, "do next func")
		}
	}
	return errors.Wrap(rows.Err(), "rows")
}

func (db *DB) QueryRow(query string, args pgx.NamedArgs, dest ...any) error {
	var row pgx.Row
	if db.tx != nil {
		row = db.tx.QueryRow(db.stopCtx, query, args)
	} else {
		row = db.pool.QueryRow(db.stopCtx, query, args)
	}
	return errors.Wrap(row.Scan(dest...), "scan")
}

func (db *DB) Ping() error {
	return db.pool.Ping(db.stopCtx)
}

func (db *DB) SendBatch(batch *pgx.Batch) (pgx.BatchResults, error) {
	if db.tx != nil {
		return db.tx.SendBatch(db.stopCtx, batch), nil
	}
	return db.pool.SendBatch(db.stopCtx, batch), nil
}

func QueryMany[T any](db *DB, query string, pointer func(*T) []any, args ...pgx.NamedArgs) ([]T, error) {
	var res = []T{}
	var arg pgx.NamedArgs
	if len(args) > 0 {
		arg = args[0]
	}

	err := db.Query(query, arg, func(rows pgx.Rows) error {
		var elem T
		if err := rows.Scan(pointer(&elem)...); err != nil {
			return errors.Wrap(err, "scan")
		}

		res = append(res, elem)
		return nil
	})

	return res, errors.Wrap(err, "query")
}

type IfaceLister interface {
	AsIfaceList() []any
}

func IfaceListFunc[T IfaceLister]() func(T) []any {
	return func(t T) []any {
		return t.AsIfaceList()
	}
}

func WrapNotFound(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return errx.ErrNotFound
	}

	return errors.Wrap(err, "query")
}

func WrapAlreadyExists(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		err = errx.ErrAlreadyExists
	}

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation {
		return errors.Wrap(errx.ErrUnprocessable, err.Error())
	}

	return errors.Wrap(err, "exec")
}

func HandleExecResult(res pgconn.CommandTag, err error) error {
	if err = WrapAlreadyExists(err); err != nil {
		return errors.Wrap(err, "exec")
	}

	if res.RowsAffected() == 0 {
		return errors.Wrap(errx.ErrNotFound, "check affected")
	}

	return nil
}
