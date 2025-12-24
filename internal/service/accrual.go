package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmnAlex/gophermart/internal/consts"
	"github.com/dmnAlex/gophermart/internal/consts/accrualstatus"
	"github.com/dmnAlex/gophermart/internal/consts/orderstatus"
	"github.com/dmnAlex/gophermart/internal/logger"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) fetchNextOrderBatch() error {
	orders, err := s.repo.LockAndGetOrderBatch(consts.OrderBatchSize)
	if err != nil {
		return errors.Wrap(err, "lock and get order")
	}

	for i := range orders {
		select {
		case <-s.ordersStopChan:
			return nil
		case s.ordersChan <- &orders[i]:
		}
	}

	return nil
}

func (s *service) processOrder(order *model.Order) error {
	defer func() {
		if err := s.repo.UpdateOrder(order.ID, order.Status, order.Accrual); err != nil {
			logger.Log.Error("update order", zap.Error(err))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), consts.AccrualTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/orders/%s", s.cfg.AccrualSystemAddress, order.Number), nil)
	if err != nil {
		return errors.Wrap(err, "create request with context")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "read body")
	}

	var payload model.AccrualResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return errors.Wrap(err, "unmarshal payload")
	}

	switch payload.Status {
	case accrualstatus.Processing:
		order.Status = orderstatus.Processing
	case accrualstatus.Invalid:
		order.Status = orderstatus.Invalid
	case accrualstatus.Processed:
		order.Status = orderstatus.Processed
		order.Accrual = payload.Accrual
	}

	return nil
}

func (s *service) freeStaleLocks() error {
	threshold := time.Now().Add(-consts.OrderLockTimeout)

	return s.repo.FreeStaleLocks(threshold)
}
