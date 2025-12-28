package server

import (
	"github.com/dmnAlex/gophermart/internal/config"
	"github.com/dmnAlex/gophermart/internal/handler"
	"github.com/dmnAlex/gophermart/internal/middleware"
	"github.com/gin-gonic/gin"
)

func newRouter(h *handler.Handler, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())

	r.POST("/api/user/register", h.HandlePostAPIUserRegister)
	r.POST("/api/user/login", h.HandlePostAPIUserLogin)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(cfg))

	auth.GET("/ping", h.HandlePing)

	auth.POST("/api/user/orders", h.HandlePostAPIUserAddOrder)
	auth.GET("/api/user/orders", h.HandleGetAPIUserGetOrders)

	auth.GET("/api/user/balance", h.HandleGetAPIUserBalance)
	auth.POST("/api/user/balance/withdraw", h.HandlePostAPIUserBalanceWithdraw)
	auth.GET("/api/user/withdrawals", h.HandleGetAPIUserWithdrawals)

	return r
}
