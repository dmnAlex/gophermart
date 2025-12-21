package consts

import "time"

const (
	AuthTokenName     = "auth_token"
	AuthTokenDutation = 1 * time.Hour

	AccrualTimeout     = 10 * time.Second
	AccrualWorkerCount = 5

	OrderLockTimeout   = 3 * time.Minute
	OrderLockFreeDelay = 1 * time.Second

	OrderChanSize  = 100
	OrderBatchSize = 10
)
