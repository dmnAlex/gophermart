package service

import (
	"context"

	"github.com/dmnAlex/gophermart/internal/config"
	"github.com/dmnAlex/gophermart/internal/consts"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	RegisterUser(login, password string) (uuid.UUID, error)
	CheckPassword(login, password string) (uuid.UUID, error)

	AddOrder(number string, userID uuid.UUID) error
	GetAllOrders(userID uuid.UUID) ([]model.Order, error)

	GetBalance(userID uuid.UUID) (model.Balance, error)
	AddWithdrawal(userID uuid.UUID, number string, sum float64) error
	GetAllWithdrawals(userID uuid.UUID) ([]model.Withdrawal, error)

	StartAccrualWorkers(ctx context.Context)
	StopAccrualWorkers() error

	Ping() error
}

type GophermartService struct {
	repo       repository.Repository
	ordersChan chan *model.Order
	cfg        *config.Config
	eg         *errgroup.Group
}

func NewGophermartService(repo repository.Repository, cfg *config.Config) *GophermartService {
	return &GophermartService{
		repo:       repo,
		ordersChan: make(chan *model.Order, consts.OrderChanSize),
		cfg:        cfg,
	}
}

func (s *GophermartService) Ping() error {
	return errors.Wrap(s.repo.Ping(), "ping")
}
