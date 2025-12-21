package repository

import (
	"time"

	"github.com/dmnAlex/gophermart/internal/consts/orderstatus"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/storage/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type RepoIface interface {
	AddUser(login, passwordHash string) (uuid.UUID, error)
	GetByLogin(login string) (uuid.UUID, string, error)

	AddOrder(number string, userID uuid.UUID) error
	GetOrderUserID(number string) (uuid.UUID, error)
	GetOrdersByLogin(userID uuid.UUID) ([]model.Order, error)

	GetBalance(userID uuid.UUID) (model.Balance, error)
	LockUserForUpdate(userID uuid.UUID) error
	AddWithdrawal(userID uuid.UUID, number string, sum float64) (uuid.UUID, error)
	GetAllWithdrawals(userID uuid.UUID) ([]model.Withdrawal, error)

	UpdateOrder(id uuid.UUID, status orderstatus.Type, accrual *float64) error
	LockAndGetOrderBatch(batchSize int) ([]model.Order, error)
	FreeStaleLocks(threshold time.Time) error

	DoTx(f func(r *Repo) error, opts ...*pgx.TxOptions) error
	Ping() error
	Close() error
}

type Repo struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) DoTx(f func(r *Repo) error, opts ...*pgx.TxOptions) error {
	return r.db.DoTx(func(db *pg.DB) error {
		return f(NewRepository(db))
	}, opts...)
}

func (r *Repo) Ping() error {
	return errors.Wrap(r.db.Ping(), "ping")
}

func (r *Repo) Close() error {
	return errors.Wrap(r.db.Close(), "close")
}
