package service

import (
	"context"
	"time"

	"github.com/dmnAlex/gophermart/internal/consts"
	"github.com/dmnAlex/gophermart/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (s *GophermartService) StartAccrualWorkers(ctx context.Context) {
	s.eg, ctx = errgroup.WithContext(ctx)

	for i := range consts.AccrualWorkerCount {
		s.eg.Go(func() error {
			return s.accrualWorker(ctx, i)
		})
	}

	s.eg.Go(func() error {
		return s.staleLocksWorker(ctx)
	})

	s.eg.Go(func() error {
		return s.fetchWorker(ctx)
	})
}

func (s *GophermartService) StopAccrualWorkers() error {
	return s.eg.Wait()
}

func (s *GophermartService) accrualWorker(ctx context.Context, workerID int) error {
	logger.Log.Info("start accrual worker", zap.Int("worker_id", workerID))

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("stop accrual worker", zap.Int("worker_id", workerID))
			return nil
		case order := <-s.ordersChan:
			if err := s.processOrder(order); err != nil {
				logger.Log.Error("process order",
					zap.Int("worker_id", workerID),
					zap.String("order", order.Number),
					zap.Error(err),
				)
			}
		}
	}
}

func (s *GophermartService) staleLocksWorker(ctx context.Context) error {
	logger.Log.Info("start stale locks worker")

	ticker := time.NewTicker(consts.LockFreeDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("stop stale locks worker")
			return nil
		case <-ticker.C:
			if err := s.freeStaleLocks(); err != nil {
				logger.Log.Error("free stale locks", zap.Error(err))
			}
		}
	}
}

func (s *GophermartService) fetchWorker(ctx context.Context) error {
	logger.Log.Info("start fetch worker")

	ticker := time.NewTicker(consts.FetchOrdersDelay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("stop fetch worker")
			return nil
		case <-ticker.C:
			if err := s.fetchNextOrderBatch(ctx); err != nil {
				logger.Log.Error("fetch next order batch", zap.Error(err))
			}
		}
	}
}
