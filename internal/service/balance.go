package service

import (
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/model/errx"
	"github.com/dmnAlex/gophermart/internal/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *service) GetBalance(userID uuid.UUID) (model.Balance, error) {
	return s.repo.GetBalance(userID)
}

func (s *service) AddWithdrawal(userID uuid.UUID, number string, sum float64) error {
	err := s.repo.DoTx(func(rTx *repository.Repo) error {
		if err := rTx.LockUserForUpdate(userID); err != nil {
			return errors.Wrap(err, "lock user for update")
		}

		_, err := rTx.AddWithdrawal(userID, number, sum)
		if errors.Is(err, errx.ErrNotFound) {
			return errx.ErrInsufficientBalance
		}

		return errors.Wrap(err, "add withdrawal")
	})

	return errors.Wrap(err, "do tx")
}

func (s *service) GetAllWithdrawals(userID uuid.UUID) ([]model.Withdrawal, error) {
	return s.repo.GetAllWithdrawals(userID)
}
