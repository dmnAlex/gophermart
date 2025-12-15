package main

import (
	"fmt"
	"net/http"

	"github.com/dmnAlex/gophermart/internal/config"
	"github.com/dmnAlex/gophermart/internal/consts"
	"github.com/dmnAlex/gophermart/internal/handler"
	"github.com/dmnAlex/gophermart/internal/logger"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(consts.AuthTokenName)
		if err != nil {
			c.String(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid || claims.Login == "" {
			c.String(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		c.Set("caller", &model.Caller{Login: claims.Login})
		c.Next()
	}
}

func newRouter(h *handler.Handler, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.LoggerMiddleware())

	r.POST("/api/user/register", h.HandleAPIUserRegister)
	r.POST("/api/user/login", h.HandleAPIUserLogin)

	auth := r.Group("/")
	auth.Use(AuthMiddleware(cfg))
	auth.GET("/ping", h.HandlePing)

	return r
}
