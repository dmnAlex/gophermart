package consts

import "time"

const (
	AuthTokenName     = "auth_token"
	AuthTokenDuration = 1 * time.Hour

	AccrualTimeout     = 10 * time.Second
	AccrualWorkerCount = 5

	OrderLockTimeout = 3 * time.Minute
	LockFreeDelay    = 1 * time.Second
	FetchOrdersDelay = 500 * time.Millisecond

	OrderChanSize  = 100
	OrderBatchSize = 10
)
