package handler

import (
	"net/http"

	"github.com/dmnAlex/gophermart/internal/config"
	"github.com/dmnAlex/gophermart/internal/logger"
	"github.com/dmnAlex/gophermart/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
	config  *config.Config
}

func NewHandler(s service.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: s,
		config:  cfg,
	}
}

func (h *Handler) HandlePing(c *gin.Context) {
	if err := h.service.Ping(); err != nil {
		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	c.Status(http.StatusOK)
}
