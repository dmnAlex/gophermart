package repository

import (
	"github.com/dmnAlex/gophermart/internal/storage/pg"
	"github.com/jackc/pgx/v5"
)

const addUserSQL = `
	INSERT INTO users (login, password)
	VALUES (@login, @password)
`

func (r *repository) AddUser(login, password string) error {
	args := pgx.NamedArgs{
		"login":    login,
		"password": password,
	}

	res, err := r.db.Exec(addUserSQL, args)

	return pg.HandleExecResult(res, err)
}

const getPasswordSQL = `
	SELECT password
	FROM users
	WHERE login = @login
`

func (r *repository) GetPassword(login string) (string, error) {
	args := pgx.NamedArgs{
		"login": login,
	}

	var password string
	err := r.db.QueryRow(getPasswordSQL, args, &password)

	return password, pg.WrapNotFound(err)
}
