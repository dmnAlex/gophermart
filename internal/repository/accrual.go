package repository

import (
	"time"

	"github.com/dmnAlex/gophermart/internal/consts/orderstatus"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/storage/pg"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

const updateOrderSQL = `
	UPDATE orders
	SET status = @status, accrual = @accrual, updated_at = NOW(), is_locked = FALSE
	WHERE id = @id
`

func (r *GophermartRepository) UpdateOrder(id uuid.UUID, status orderstatus.Type, accrual *float64) error {
	args := pgx.NamedArgs{
		"id":      id,
		"status":  status,
		"accrual": accrual,
	}

	res, err := r.db.Exec(updateOrderSQL, args)

	return pg.HandleExecResult(res, err)
}

const lockAndGetOrderBatchSQL = `
	WITH selected_orders AS (
    SELECT id
    FROM orders
    WHERE status IN ('NEW', 'PROCESSING')
      AND is_locked = FALSE
    ORDER BY updated_at ASC
    LIMIT @batch_size
    FOR UPDATE SKIP LOCKED
	)
	UPDATE orders o
	SET is_locked  = TRUE,
		updated_at = NOW()
	FROM selected_orders so
	WHERE o.id = so.id
	RETURNING 
		o.id, o.number, o.status, o.accrual, o.uploaded_at
`

func (r *GophermartRepository) LockAndGetOrderBatch(batchSize int) ([]model.Order, error) {
	args := pgx.NamedArgs{
		"batch_size": batchSize,
	}
	return pg.QueryMany(r.db, lockAndGetOrderBatchSQL, pg.IfaceListFunc[*model.Order](), args)
}

const freeStaleLocksSQL = `
	UPDATE orders
	SET is_locked = FALSE
	WHERE is_locked = TRUE 
		AND updated_at <= @threshold
`

func (r *GophermartRepository) FreeStaleLocks(threshold time.Time) error {
	args := pgx.NamedArgs{
		"threshold": threshold,
	}

	_, err := r.db.Exec(freeStaleLocksSQL, args)

	return errors.Wrap(err, "exec")
}
