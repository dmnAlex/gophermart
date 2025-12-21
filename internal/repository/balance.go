package repository

import (
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/storage/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const getBalanceSQL = `
	SELECT
		COALESCE(SUM(o.accrual), 0) - COALESCE(SUM(w.sum), 0) as current,
		COALESCE(SUM(w.sum), 0) as withdrawn
	FROM users u
	LEFT JOIN orders o ON o.user_id = u.id AND o.status = 'PROCESSED'
	LEFT JOIN withdrawals w ON w.user_id = u.id
	WHERE u.id = @user_id
	GROUP BY u.id
`

func (r *Repo) GetBalance(userID uuid.UUID) (model.Balance, error) {
	args := pgx.NamedArgs{
		"user_id": userID,
	}

	var balance model.Balance
	err := r.db.QueryRow(getBalanceSQL, args, balance.AsIfaceList()...)

	return balance, pg.WrapNotFound(err)
}

const lockUserForUpdateSQL = `
	SELECT 1
	FROM users
	WHERE id = @id
	FOR UPDATE
`

func (r *Repo) LockUserForUpdate(userID uuid.UUID) error {
	args := pgx.NamedArgs{
		"id": userID,
	}
	_, err := r.db.Exec(lockUserForUpdateSQL, args)

	return errors.Wrap(err, "exec")
}

const addWithdrawalSQL = `
	WITH user_balance AS (
		SELECT
			(
				SELECT COALESCE(SUM(accrual), 0) 
				FROM orders 
				WHERE user_id = @user_id 
					AND status = 'PROCESSED'
			) - (
				SELECT COALESCE(SUM(sum), 0) 
				FROM withdrawals 
				WHERE user_id = @user_id
			) as available
	),
	attempt_insert AS (
		INSERT INTO withdrawals (user_id, number, sum)
		SELECT @user_id, @number, @sum
		FROM user_balance
		WHERE user_balance.available >= @sum
		RETURNING id
	)
	SELECT id FROM attempt_insert
`

func (r *Repo) AddWithdrawal(userID uuid.UUID, number string, sum float64) (uuid.UUID, error) {
	args := pgx.NamedArgs{
		"user_id": userID,
		"number":  number,
		"sum":     sum,
	}

	var withdrawalID uuid.UUID
	err := r.db.QueryRow(addWithdrawalSQL, args, &withdrawalID)

	return withdrawalID, pg.WrapNotFound(err)
}

const getWithdrawalsSQL = `
	SELECT number, sum, processed_at
	FROM withdrawals
	WHERE user_id = @user_id
`

func (r *Repo) GetAllWithdrawals(userID uuid.UUID) ([]model.Withdrawal, error) {
	args := pgx.NamedArgs{
		"user_id": userID,
	}

	return pg.QueryMany(r.db, string(getWithdrawalsSQL), pg.IfaceListFunc[*model.Withdrawal](), args)
}
