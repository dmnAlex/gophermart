package handler

import (
	"net/http"
	"time"

	"github.com/dmnAlex/gophermart/internal/logger"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/model/errx"
	"github.com/dmnAlex/gophermart/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *Handler) HandlePostAPIUserAddOrder(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil || len(body) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	number := string(body)
	if !utils.IsValidLuhn(number) {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	caller := c.MustGet("caller").(*model.Caller)
	if err := h.service.AddOrder(number, caller.UserID); err != nil {
		switch {
		case errors.Is(err, errx.ErrAlreadyAccepted):
			c.Status(http.StatusOK)
		case errors.Is(err, errx.ErrConflict):
			c.Status(http.StatusConflict)
		default:
			c.Status(http.StatusInternalServerError)
			err = errors.Wrap(err, "add order")
			logger.Log.Error(err.Error())
		}
		return
	}

	c.Status(http.StatusAccepted)
}

func (h *Handler) HandleGetAPIUserGetOrders(c *gin.Context) {
	caller := c.MustGet("caller").(*model.Caller)
	orders, err := h.service.GetAllOrders(caller.UserID)
	if err != nil {
		err = errors.Wrap(err, "get all orders")
		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	for i := range orders {
		orders[i].UploadedAt = orders[i].UploadedAt.Truncate(time.Second)
	}

	c.JSON(http.StatusOK, orders)
}
