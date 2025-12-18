package service

import (
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type ServiceIface interface {
	RegisterUser(login, password string) (uuid.UUID, error)
	CheckPassword(login, password string) (uuid.UUID, error)

	AddOrder(number string, userID uuid.UUID) error
	GetAllOrders(userID uuid.UUID) ([]model.Order, error)

	GetBalance(userID uuid.UUID) (model.Balance, error)
	AddWithdrawal(userID uuid.UUID, number string, sum decimal.Decimal) error
	GetAllWithdrawals(userID uuid.UUID) ([]model.Withdrawal, error)

	Ping() error
}

type service struct {
	repo repository.RepoIface
}

func NewService(repo repository.RepoIface) ServiceIface {
	return &service{repo: repo}
}

func (s *service) Ping() error {
	return errors.Wrap(s.repo.Ping(), "ping")
}
