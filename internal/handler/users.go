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
	"github.com/pkg/errors"
)

func (h *Handler) HandleAPIUserRegister(c *gin.Context) {
	var req model.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, errx.ErrBadRequest.Error())
		return
	}

	if err := h.service.RegisterUser(req.Login, req.Password); err != nil {
		err = errors.Wrap(err, "register user")
		if errors.Is(err, errx.ErrAlreadyExists) {
			c.String(http.StatusConflict, errx.ErrAlreadyExists.Error())
			return
		}

		c.String(http.StatusInternalServerError, errx.ErrInternalError.Error())
		logger.Log.Error(err.Error())
		return
	}

	if err := h.setCookie(c, req.Login); err != nil {
		err = errors.Wrap(err, "set cookie")
		c.String(http.StatusInternalServerError, errx.ErrInternalError.Error())
		logger.Log.Error(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) HandleAPIUserLogin(c *gin.Context) {
	var req model.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, errx.ErrBadRequest.Error())
		return
	}

	if err := h.service.CheckPassword(req.Login, req.Password); err != nil {
		if errors.Is(err, errx.ErrUnauthorized) {
			c.String(http.StatusUnauthorized, errx.ErrUnauthorized.Error())
			return
		}

		c.String(http.StatusInternalServerError, errx.ErrInternalError.Error())
		logger.Log.Error(err.Error())
		return
	}

	if err := h.setCookie(c, req.Login); err != nil {
		err = errors.Wrap(err, "set cookie")
		c.String(http.StatusInternalServerError, errx.ErrInternalError.Error())
		logger.Log.Error(err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) setCookie(c *gin.Context, login string) error {
	claims := &model.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.AuthTokenDutation)),
		},
		Login: login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		return errors.Wrap(err, "create token with claims")
	}
	c.SetCookie(consts.AuthTokenName, signedToken, int(consts.AuthTokenDutation.Seconds()), "/", "", false, true)

	return nil
}
