package service

import (
	"github.com/dmnAlex/gophermart/internal/config"
	"github.com/dmnAlex/gophermart/internal/consts"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ServiceIface interface {
	RegisterUser(login, password string) (uuid.UUID, error)
	CheckPassword(login, password string) (uuid.UUID, error)

	AddOrder(number string, userID uuid.UUID) error
	GetAllOrders(userID uuid.UUID) ([]model.Order, error)

	GetBalance(userID uuid.UUID) (model.Balance, error)
	AddWithdrawal(userID uuid.UUID, number string, sum float64) error
	GetAllWithdrawals(userID uuid.UUID) ([]model.Withdrawal, error)

	FreeStaleLocks() error

	Ping() error
}

type service struct {
	repo       repository.RepoIface
	ordersChan chan *model.Order
	cfg        *config.Config
}

func NewService(repo repository.RepoIface, cfg *config.Config) ServiceIface {
	return &service{
		repo:       repo,
		ordersChan: make(chan *model.Order, consts.OrderChanSize),
		cfg:        cfg,
	}
}

func (s *service) Ping() error {
	return errors.Wrap(s.repo.Ping(), "ping")
}
