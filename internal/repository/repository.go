package repository

import (
	"github.com/dmnAlex/gophermart/internal/storage/pg"
	"github.com/pkg/errors"
)

type RepoIface interface {
	AddUser(login, password string) error
	GetPassword(login string) (string, error)

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
