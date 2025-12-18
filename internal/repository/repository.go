package repository

import (
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/storage/pg"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type RepoIface interface {
	AddUser(login, passwordHash string) (uuid.UUID, error)
	GetByLogin(login string) (uuid.UUID, string, error)

	AddOrder(number string, userID uuid.UUID) error
	GetOrderUserID(number string) (uuid.UUID, error)
	GetOrdersByLogin(userID uuid.UUID) ([]model.Order, error)

	GetBalance(userID uuid.UUID) (model.Balance, error)
	AddWithdrawal(userID uuid.UUID, number string, sum decimal.Decimal) (uuid.UUID, error)
	GetAllWithdrawals(userID uuid.UUID) ([]model.Withdrawal, error)

	Ping() error
	Close() error
}

type repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) RepoIface {
	return &repository{db: db}
}

func (r *repository) Ping() error {
	return errors.Wrap(r.db.Ping(), "ping")
}

func (r *repository) Close() error {
	return errors.Wrap(r.db.Close(), "close")
}
