package service

import (
	"github.com/dmnAlex/gophermart/internal/model/errx"
	"github.com/dmnAlex/gophermart/internal/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *GophermartService) RegisterUser(login, password string) (uuid.UUID, error) {
	passwordHash := utils.Sha256Hex([]byte(password))

	return s.repo.AddUser(login, passwordHash)
}

func (s *GophermartService) CheckPassword(login, password string) (uuid.UUID, error) {
	passwordHash := utils.Sha256Hex([]byte(password))

	userID, storedPasswordHash, err := s.repo.GetByLogin(login)
	if err != nil {
		if errors.Is(err, errx.ErrNotFound) {
			return uuid.Nil, errx.ErrUnauthorized
		}

		return uuid.Nil, errors.Wrap(err, "get password")
	}

	if passwordHash != storedPasswordHash {
		return uuid.Nil, errx.ErrUnauthorized
	}

	return userID, nil
}
