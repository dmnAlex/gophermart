package service

import (
	"github.com/dmnAlex/gophermart/internal/model/errx"
	"github.com/dmnAlex/gophermart/internal/utils"
	"github.com/pkg/errors"
)

func (s *service) RegisterUser(login, password string) error {
	hexPassword := utils.Sha256Hex([]byte(password))

	return s.repo.AddUser(login, hexPassword)
}

func (s *service) CheckPassword(login, password string) error {
	hexPassword := utils.Sha256Hex([]byte(password))

	password, err := s.repo.GetPassword(login)
	if err != nil {
		if errors.Is(err, errx.ErrNotFound) {
			return errx.ErrUnauthorized
		}

		return errors.Wrap(err, "get password")
	}

	if hexPassword != password {
		return errx.ErrUnauthorized
	}

	return nil
}
