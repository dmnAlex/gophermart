package model

import (
	"time"

	"github.com/dmnAlex/gophermart/internal/consts/accrualstatus"
	"github.com/dmnAlex/gophermart/internal/consts/orderstatus"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
	ID         uuid.UUID        `json:"-"`
	Number     string           `json:"number"`
	Status     orderstatus.Type `json:"status"`
	Accrual    *float64         `json:"accrual,omitempty"`
	UploadedAt time.Time        `json:"uploaded_at"`
}

func (m *Order) AsIfaceList() []any {
	return []any{&m.ID, &m.Number, &m.Status, &m.Accrual, &m.UploadedAt}
}

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (m *Withdrawal) AsIfaceList() []any {
	return []any{&m.Order, &m.Sum, &m.ProcessedAt}
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (m *Balance) AsIfaceList() []any {
	return []any{&m.Current, &m.Withdrawn}
}

type WithdrawalRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type AccrualResponse struct {
	Order   string             `json:"order"`
	Status  accrualstatus.Type `json:"status"`
	Accrual *float64           `json:"accrual"`
}
