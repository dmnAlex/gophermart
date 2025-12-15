package model

import (
	"time"

	"github.com/dmnAlex/gophermart/internal/consts/orderstatus"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Caller struct {
	UserID uuid.UUID
}

type AuthRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

type Order struct {
	Number     string           `json:"number"`
	Status     orderstatus.Type `json:"status"`
	Accrual    *decimal.Decimal `json:"accrual,omitempty"`
	UploadedAt time.Time        `json:"uploaded_at"`
}

func (m *Order) AsIfaceList() []any {
	return []any{&m.Number, &m.Status, &m.Accrual, &m.UploadedAt}
}
