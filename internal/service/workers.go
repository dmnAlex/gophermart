package service

import (
	"time"

	"github.com/dmnAlex/gophermart/internal/consts"
	"github.com/dmnAlex/gophermart/internal/logger"
	"go.uber.org/zap"
)

func (s *service) StartAccrualWorkers() {
	for i := range consts.AccrualWorkerCount {
		s.wg.Add(1)
		go s.accrualWorker(i)
	}

	s.wg.Add(1)
	go s.staleLocksWorker()

	s.wg.Add(1)
	go s.fetchWorker()
}

func (s *service) StopAccrualWorkers() {
	close(s.ordersStopChan)
	s.wg.Wait()
}

func (s *service) accrualWorker(workerID int) {
	defer s.wg.Done()

	logger.Log.Info("start accrual worker", zap.Int("worker_id", workerID))

	for {
		select {
		case <-s.ordersStopChan:
			logger.Log.Info("stop accrual worker", zap.Int("worker_id", workerID))
			return
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

func (s *service) staleLocksWorker() {
	defer s.wg.Done()

	logger.Log.Info("start stale locks worker")

	ticker := time.NewTicker(consts.LockFreeDelay)
	defer ticker.Stop()

	for {
		select {
		case <-s.ordersStopChan:
			logger.Log.Info("stop stale locks worker")
			return
		case <-ticker.C:
			if err := s.freeStaleLocks(); err != nil {
				logger.Log.Error("free stale locks", zap.Error(err))
			}
		}
	}
}

func (s *service) fetchWorker() {
	defer s.wg.Done()

	logger.Log.Info("start fetch worker")

	ticker := time.NewTicker(consts.FetchOrdersDelay)
	defer ticker.Stop()

	for {
		select {
		case <-s.ordersStopChan:
			logger.Log.Info("stop fetch worker")
			return
		case <-ticker.C:
			if err := s.fetchNextOrderBatch(); err != nil {
				logger.Log.Error("fetch next order batch", zap.Error(err))
			}
		}
	}
}
