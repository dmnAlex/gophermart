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

func (h *Handler) HandleGetAPIUserBalance(c *gin.Context) {
	caller := c.MustGet("caller").(*model.Caller)

	balance, err := h.service.GetBalance(caller.UserID)
	if err != nil {
		err = errors.Wrap(err, "get balance")
		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (h *Handler) HandlePostAPIUserBalanceWithdraw(c *gin.Context) {
	var req model.WithdrawalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if !utils.IsValidLuhn(req.Order) {
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	caller := c.MustGet("caller").(*model.Caller)
	if err := h.service.AddWithdrawal(caller.UserID, req.Order, req.Sum); err != nil {
		if errors.Is(err, errx.ErrInsufficientBalance) {
			c.Status(http.StatusPaymentRequired)
			return
		}

		err = errors.Wrap(err, "add withdrawal")
		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) HandleGetAPIUserWithdrawals(c *gin.Context) {
	caller := c.MustGet("caller").(*model.Caller)
	withdrawals, err := h.service.GetAllWithdrawals(caller.UserID)
	if err != nil {
		err = errors.Wrap(err, "get all withdrawals")
		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	for i := range withdrawals {
		withdrawals[i].ProcessedAt = withdrawals[i].ProcessedAt.Truncate(time.Second)
	}

	c.JSON(http.StatusOK, withdrawals)
}
