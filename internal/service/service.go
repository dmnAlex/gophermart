package service

import (
	"github.com/dmnAlex/gophermart/internal/repository"
	"github.com/pkg/errors"
)

type ServiceIface interface {
	RegisterUser(login, password string) error
	CheckPassword(login, password string) error

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
