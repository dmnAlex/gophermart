package handler

import (
	"net/http"
	"time"

	"github.com/dmnAlex/gophermart/internal/consts"
	"github.com/dmnAlex/gophermart/internal/logger"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/model/errx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *Handler) HandlePostAPIUserRegister(c *gin.Context) {
	var req model.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userID, err := h.service.RegisterUser(req.Login, req.Password)
	if err != nil {
		err = errors.Wrap(err, "register user")
		if errors.Is(err, errx.ErrAlreadyExists) {
			c.Status(http.StatusConflict)
			return
		}

		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	if err := h.setCookie(c, userID); err != nil {
		err = errors.Wrap(err, "set cookie")
		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) HandlePostAPIUserLogin(c *gin.Context) {
	var req model.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userID, err := h.service.CheckPassword(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, errx.ErrUnauthorized) {
			c.Status(http.StatusUnauthorized)
			return
		}

		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	if err := h.setCookie(c, userID); err != nil {
		err = errors.Wrap(err, "set cookie")
		c.Status(http.StatusInternalServerError)
		logger.Log.Error(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) setCookie(c *gin.Context, userID uuid.UUID) error {
	claims := &model.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.AuthTokenDutation)),
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		return errors.Wrap(err, "create token with claims")
	}
	c.SetCookie(consts.AuthTokenName, signedToken, int(consts.AuthTokenDutation.Seconds()), "/", "", false, true)

	return nil
}
