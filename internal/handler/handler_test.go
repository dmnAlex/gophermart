package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmnAlex/gophermart/internal/config"
	"github.com/dmnAlex/gophermart/internal/consts"
	"github.com/dmnAlex/gophermart/internal/model"
	"github.com/dmnAlex/gophermart/internal/model/errx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) RegisterUser(login, password string) (uuid.UUID, error) {
	args := m.Called(login, password)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockService) CheckPassword(login, password string) (uuid.UUID, error) {
	args := m.Called(login, password)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockService) AddOrder(number string, userID uuid.UUID) error {
	args := m.Called(number, userID)
	return args.Error(0)
}

func (m *MockService) GetAllOrders(userID uuid.UUID) ([]model.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockService) GetBalance(userID uuid.UUID) (model.Balance, error) {
	args := m.Called(userID)
	return args.Get(0).(model.Balance), args.Error(1)
}

func (m *MockService) AddWithdrawal(userID uuid.UUID, number string, sum float64) error {
	args := m.Called(userID, number, sum)
	return args.Error(0)
}

func (m *MockService) GetAllWithdrawals(userID uuid.UUID) ([]model.Withdrawal, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Withdrawal), args.Error(1)
}

func (m *MockService) StartAccrualWorkers(ctx context.Context) {
	m.Called()
}

func (m *MockService) StopAccrualWorkers() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockService) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func setupTest() (*gin.Engine, *MockService, *Handler) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockService)
	cfg := &config.Config{
		JWTSecret: "test-secret",
	}
	handler := NewHandler(mockService, cfg)

	router := gin.New()
	return router, mockService, handler
}

// Tests for HandlePostAPIUserRegister
func TestHandlePostAPIUserRegister_Success(t *testing.T) {
	router, mockService, handler := setupTest()
	router.POST("/register", handler.HandlePostAPIUserRegister)

	userID := uuid.New()
	authReq := model.AuthRequest{
		Login:    "testuser",
		Password: "password123",
	}

	mockService.On("RegisterUser", authReq.Login, authReq.Password).Return(userID, nil)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	cookies := res.Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, consts.AuthTokenName, cookies[0].Name)
	assert.NotEmpty(t, cookies[0].Value)

	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserRegister_InvalidJSON(t *testing.T) {
	router, mockService, handler := setupTest()
	router.POST("/register", handler.HandlePostAPIUserRegister)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "RegisterUser")
}

func TestHandlePostAPIUserRegister_UserAlreadyExists(t *testing.T) {
	router, mockService, handler := setupTest()
	router.POST("/register", handler.HandlePostAPIUserRegister)

	authReq := model.AuthRequest{
		Login:    "existinguser",
		Password: "password123",
	}

	mockService.On("RegisterUser", authReq.Login, authReq.Password).Return(uuid.Nil, errx.ErrAlreadyExists)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserRegister_InternalError(t *testing.T) {
	router, mockService, handler := setupTest()
	router.POST("/register", handler.HandlePostAPIUserRegister)

	authReq := model.AuthRequest{
		Login:    "testuser",
		Password: "password123",
	}

	mockService.On("RegisterUser", authReq.Login, authReq.Password).Return(uuid.Nil, errors.New("db error"))

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for HandlePostAPIUserLogin
func TestHandlePostAPIUserLogin_Success(t *testing.T) {
	router, mockService, handler := setupTest()
	router.POST("/login", handler.HandlePostAPIUserLogin)

	userID := uuid.New()
	authReq := model.AuthRequest{
		Login:    "testuser",
		Password: "password123",
	}

	mockService.On("CheckPassword", authReq.Login, authReq.Password).Return(userID, nil)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	res := w.Result()
	defer res.Body.Close()
	cookies := res.Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, consts.AuthTokenName, cookies[0].Name)

	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserLogin_InvalidJSON(t *testing.T) {
	router, mockService, handler := setupTest()
	router.POST("/login", handler.HandlePostAPIUserLogin)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "CheckPassword")
}

func TestHandlePostAPIUserLogin_Unauthorized(t *testing.T) {
	router, mockService, handler := setupTest()
	router.POST("/login", handler.HandlePostAPIUserLogin)

	authReq := model.AuthRequest{
		Login:    "testuser",
		Password: "wrongpassword",
	}

	mockService.On("CheckPassword", authReq.Login, authReq.Password).Return(uuid.Nil, errx.ErrUnauthorized)

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserLogin_InternalError(t *testing.T) {
	router, mockService, handler := setupTest()
	router.POST("/login", handler.HandlePostAPIUserLogin)

	authReq := model.AuthRequest{
		Login:    "testuser",
		Password: "password123",
	}

	mockService.On("CheckPassword", authReq.Login, authReq.Password).Return(uuid.Nil, errors.New("db error"))

	body, _ := json.Marshal(authReq)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for HandlePostAPIUserAddOrder
func TestHandlePostAPIUserAddOrder_Success(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/order", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserAddOrder(c)
	})

	orderNumber := "4539319503436467" // Valid Luhn number
	mockService.On("AddOrder", orderNumber, userID).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(orderNumber))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserAddOrder_EmptyBody(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/order", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserAddOrder(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(""))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "AddOrder")
}

func TestHandlePostAPIUserAddOrder_InvalidLuhn(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/order", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserAddOrder(c)
	})

	orderNumber := "1234567890" // Invalid Luhn
	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(orderNumber))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	mockService.AssertNotCalled(t, "AddOrder")
}

func TestHandlePostAPIUserAddOrder_AlreadyAccepted(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/order", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserAddOrder(c)
	})

	orderNumber := "4539319503436467"
	mockService.On("AddOrder", orderNumber, userID).Return(errx.ErrAlreadyAccepted)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(orderNumber))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserAddOrder_Conflict(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/order", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserAddOrder(c)
	})

	orderNumber := "4539319503436467"
	mockService.On("AddOrder", orderNumber, userID).Return(errx.ErrConflict)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(orderNumber))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserAddOrder_InternalError(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/order", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserAddOrder(c)
	})

	orderNumber := "4539319503436467"
	mockService.On("AddOrder", orderNumber, userID).Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(orderNumber))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for HandleGetAPIUserGetOrders
func TestHandleGetAPIUserGetOrders_Success(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.GET("/orders", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandleGetAPIUserGetOrders(c)
	})

	accrual := 500.0
	orders := []model.Order{
		{
			ID:         uuid.New(),
			Number:     "4539319503436467",
			Status:     "PROCESSED",
			Accrual:    &accrual,
			UploadedAt: time.Now(),
		},
	}

	mockService.On("GetAllOrders", userID).Return(orders, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.Order
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, orders[0].Number, response[0].Number)

	mockService.AssertExpectations(t)
}

func TestHandleGetAPIUserGetOrders_NoOrders(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.GET("/orders", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandleGetAPIUserGetOrders(c)
	})

	mockService.On("GetAllOrders", userID).Return([]model.Order{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandleGetAPIUserGetOrders_InternalError(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.GET("/orders", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandleGetAPIUserGetOrders(c)
	})

	mockService.On("GetAllOrders", userID).Return([]model.Order{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for HandleGetAPIUserBalance
func TestHandleGetAPIUserBalance_Success(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.GET("/balance", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandleGetAPIUserBalance(c)
	})

	balance := model.Balance{
		Current:   1000.50,
		Withdrawn: 200.00,
	}

	mockService.On("GetBalance", userID).Return(balance, nil)

	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Balance
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, balance.Current, response.Current)
	assert.Equal(t, balance.Withdrawn, response.Withdrawn)

	mockService.AssertExpectations(t)
}

func TestHandleGetAPIUserBalance_InternalError(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.GET("/balance", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandleGetAPIUserBalance(c)
	})

	mockService.On("GetBalance", userID).Return(model.Balance{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for HandlePostAPIUserBalanceWithdraw
func TestHandlePostAPIUserBalanceWithdraw_Success(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/withdraw", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserBalanceWithdraw(c)
	})

	withdrawalReq := model.WithdrawalRequest{
		Order: "4539319503436467",
		Sum:   100.50,
	}

	mockService.On("AddWithdrawal", userID, withdrawalReq.Order, withdrawalReq.Sum).Return(nil)

	body, _ := json.Marshal(withdrawalReq)
	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserBalanceWithdraw_InvalidJSON(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/withdraw", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserBalanceWithdraw(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "AddWithdrawal")
}

func TestHandlePostAPIUserBalanceWithdraw_InvalidLuhn(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/withdraw", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserBalanceWithdraw(c)
	})

	withdrawalReq := model.WithdrawalRequest{
		Order: "1234567890", // Invalid Luhn
		Sum:   100.50,
	}

	body, _ := json.Marshal(withdrawalReq)
	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	mockService.AssertNotCalled(t, "AddWithdrawal")
}

func TestHandlePostAPIUserBalanceWithdraw_InsufficientBalance(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/withdraw", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserBalanceWithdraw(c)
	})

	withdrawalReq := model.WithdrawalRequest{
		Order: "4539319503436467",
		Sum:   1000000.00,
	}

	mockService.On("AddWithdrawal", userID, withdrawalReq.Order, withdrawalReq.Sum).Return(errx.ErrInsufficientBalance)

	body, _ := json.Marshal(withdrawalReq)
	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusPaymentRequired, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlePostAPIUserBalanceWithdraw_InternalError(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.POST("/withdraw", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandlePostAPIUserBalanceWithdraw(c)
	})

	withdrawalReq := model.WithdrawalRequest{
		Order: "4539319503436467",
		Sum:   100.50,
	}

	mockService.On("AddWithdrawal", userID, withdrawalReq.Order, withdrawalReq.Sum).Return(errors.New("db error"))

	body, _ := json.Marshal(withdrawalReq)
	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// Tests for HandleGetAPIUserWithdrawals
func TestHandleGetAPIUserWithdrawals_Success(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.GET("/withdrawals", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandleGetAPIUserWithdrawals(c)
	})

	withdrawals := []model.Withdrawal{
		{
			Order:       "4539319503436467",
			Sum:         100.50,
			ProcessedAt: time.Now(),
		},
	}

	mockService.On("GetAllWithdrawals", userID).Return(withdrawals, nil)

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []model.Withdrawal
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, withdrawals[0].Order, response[0].Order)

	mockService.AssertExpectations(t)
}

func TestHandleGetAPIUserWithdrawals_NoWithdrawals(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.GET("/withdrawals", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandleGetAPIUserWithdrawals(c)
	})

	mockService.On("GetAllWithdrawals", userID).Return([]model.Withdrawal{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandleGetAPIUserWithdrawals_InternalError(t *testing.T) {
	router, mockService, handler := setupTest()

	userID := uuid.New()
	router.GET("/withdrawals", func(c *gin.Context) {
		c.Set("caller", &model.Caller{UserID: userID})
		handler.HandleGetAPIUserWithdrawals(c)
	})

	mockService.On("GetAllWithdrawals", userID).Return([]model.Withdrawal{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
