package repository

import (
	"github.com/dmnAlex/gophermart/internal/storage/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const addUserSQL = `
	INSERT INTO users (login, password_hash)
	VALUES (@login, @password_hash)
	RETURNING id
`

func (r *Repo) AddUser(login, passwordHash string) (uuid.UUID, error) {
	args := pgx.NamedArgs{
		"login":         login,
		"password_hash": passwordHash,
	}

	var id uuid.UUID
	err := r.db.QueryRow(addUserSQL, args, &id)

	return id, pg.WrapAlreadyExists(err)
}

const getPasswordSQL = `
	SELECT id, password_hash
	FROM users
	WHERE login = @login
`

func (r *Repo) GetByLogin(login string) (uuid.UUID, string, error) {
	args := pgx.NamedArgs{
		"login": login,
	}

	var id uuid.UUID
	var password string
	err := r.db.QueryRow(getPasswordSQL, args, &id, &password)

	return id, password, pg.WrapNotFound(err)
}
