package service

import (
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/model/errx"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *GophermartService) AddOrder(number string, userID uuid.UUID) error {
	if err := s.repo.AddOrder(number, userID); err != nil {
		if errors.Is(err, errx.ErrAlreadyExists) {
			creatorUserID, err := s.repo.GetOrderUserID(number)
			if err != nil {
				return errors.Wrap(err, "get order creator login")
			}

			if creatorUserID != userID {
				return errx.ErrConflict
			}

			return errx.ErrAlreadyAccepted
		}

		return errors.Wrap(err, "add order")
	}

	return nil
}

func (s *GophermartService) GetAllOrders(userID uuid.UUID) ([]model.Order, error) {
	return s.repo.GetOrdersByLogin(userID)
}
