package repository

import (
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/storage/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const addOrderSQL = `
	INSERT INTO orders (number, user_id)
	VALUES (@number, @user_id)
`

func (r *Repo) AddOrder(number string, userID uuid.UUID) error {
	args := pgx.NamedArgs{
		"number":  number,
		"user_id": userID,
	}

	res, err := r.db.Exec(addOrderSQL, args)

	return pg.HandleExecResult(res, err)
}

const getOrderCreatorLoginSQL = `
	SELECT user_id
	FROM orders
	WHERE number = @number
`

func (r *Repo) GetOrderUserID(number string) (uuid.UUID, error) {
	var userID uuid.UUID
	args := pgx.NamedArgs{
		"number": number,
	}

	err := r.db.QueryRow(getOrderCreatorLoginSQL, args, &userID)

	return userID, pg.WrapNotFound(err)
}

const getOrdersByLogin = `
	SELECT id, number, status, accrual, uploaded_at
	FROM orders
	WHERE user_id = @user_id
	ORDER BY uploaded_at DESC
`

func (r *Repo) GetOrdersByLogin(userID uuid.UUID) ([]model.Order, error) {
	args := pgx.NamedArgs{
		"user_id": userID,
	}

	return pg.QueryMany(r.db, getOrdersByLogin, pg.IfaceListFunc[*model.Order](), args)
}
