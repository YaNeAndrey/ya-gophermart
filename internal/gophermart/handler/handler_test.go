package handler

import (
	"context"
	"errors"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/consterror"
	status "github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/status"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/mocks"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBalanceGET(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	type args struct {
		req *http.Request
		st  *mocks.MockStorageRepo
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			name: "First test",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/api/user/balance", nil),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Second test",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/api/user/balance", nil),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusOK,
		},
	}

	tests[0].args.st.EXPECT().GetUserBalance(ctx, "1").Return(nil, errors.New(""))
	tests[1].args.st.EXPECT().GetUserBalance(ctx, "1").Return(&storage.Balance{
		Current:   1.1,
		Withdrawn: 1.2,
	}, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IjEifQ.89buB88ZZu4dHpaH7b229SHhQe67gq-5Pgig2xiKm48",
			})
			w := httptest.NewRecorder()
			r := chi.NewRouter()
			st := storage.StorageRepo(tt.args.st)
			r.Get("/api/user/balance", func(rw http.ResponseWriter, r *http.Request) {
				BalanceGET(rw, r, &st)
			})
			r.ServeHTTP(w, tt.args.req)

			result := w.Result()
			require.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestBalanceWithdrawPOST(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	type args struct {
		req *http.Request
		st  *mocks.MockStorageRepo
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			name: "First test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader("{\"order\": \"1\", \"sum\":248.08}")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Second test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader("{\"order\": \"2\", \"sum\":2480000.08}")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusPaymentRequired,
		},
		{
			name: "Third test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader("{\"order\": \"3551\", \"sum\":2480000.08}")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusUnprocessableEntity,
		},
	}

	tests[0].args.st.EXPECT().DoRebiting(ctx, "1", "1", 248.08).Return(nil)
	tests[1].args.st.EXPECT().DoRebiting(ctx, "1", "2", 2480000.08).Return(consterror.ErrInsufficientFunds)
	tests[2].args.st.EXPECT().DoRebiting(ctx, "1", "3551", 2480000.08).Return(consterror.ErrOrderNotFound)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IjEifQ.89buB88ZZu4dHpaH7b229SHhQe67gq-5Pgig2xiKm48",
			})
			w := httptest.NewRecorder()
			r := chi.NewRouter()

			st := storage.StorageRepo(tt.args.st)
			r.Post("/api/user/balance/withdraw", func(rw http.ResponseWriter, r *http.Request) {
				BalanceWithdrawPOST(rw, r, &st)
			})
			r.ServeHTTP(w, tt.args.req)

			result := w.Result()
			require.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestLoginPOST(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	type args struct {
		req *http.Request
		st  *mocks.MockStorageRepo
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			name: "First test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader("{\"login\": \"1\",\"password\": \"1\"}")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Second test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader("{\"login\": \"1\",\"password\": \"2\"}")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "Third test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader("{\"login\": \"1\",\"password\": \"1\"}")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusUnauthorized,
		},
	}

	tests[0].args.st.EXPECT().CheckUserPassword(ctx, "1", "1").Return(true, nil)
	tests[1].args.st.EXPECT().CheckUserPassword(ctx, "1", "2").Return(false, nil)
	tests[2].args.st.EXPECT().CheckUserPassword(ctx, "1", "1").Return(false, consterror.ErrLoginNotFound)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IjEifQ.89buB88ZZu4dHpaH7b229SHhQe67gq-5Pgig2xiKm48",
			})
			w := httptest.NewRecorder()
			r := chi.NewRouter()

			st := storage.StorageRepo(tt.args.st)
			r.Post("/api/user/login", func(rw http.ResponseWriter, r *http.Request) {
				LoginPOST(rw, r, &st)
			})
			r.ServeHTTP(w, tt.args.req)

			result := w.Result()
			require.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestOrdersPOST(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	type args struct {
		req *http.Request
		st  *mocks.MockStorageRepo
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			name: "First test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusAccepted,
		},
		{
			name: "Second test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Second test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678903")),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusConflict,
		},
	}

	tests[0].args.st.EXPECT().AddNewOrder(ctx, "1", "12345678903").Return(&storage.Order{
		Number:     "12345678903",
		Status:     status.New,
		Accrual:    0,
		UploadDate: time.Now(),
		Sum:        0,
	}, nil)
	tests[1].args.st.EXPECT().AddNewOrder(ctx, "1", "12345678903").Return(nil, consterror.ErrDuplicateUserOrder)
	tests[2].args.st.EXPECT().AddNewOrder(ctx, "1", "12345678903").Return(nil, consterror.ErrDuplicateAnotherUserOrder)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IjEifQ.89buB88ZZu4dHpaH7b229SHhQe67gq-5Pgig2xiKm48",
			})
			w := httptest.NewRecorder()
			r := chi.NewRouter()

			st := storage.StorageRepo(tt.args.st)
			r.Post("/api/user/orders", func(rw http.ResponseWriter, r *http.Request) {
				OrdersPOST(rw, r, &st)
			})
			r.ServeHTTP(w, tt.args.req)

			result := w.Result()
			require.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestWithdrawalsGET(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	type args struct {
		req *http.Request
		st  *mocks.MockStorageRepo
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			name: "First test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/withdrawals", nil),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Second test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/withdrawals", nil),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Third test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/withdrawals", nil),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Fourth test",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/api/user/withdrawals", nil),
				st:  mocks.NewMockStorageRepo(ctrl),
			},
			statusCode: http.StatusUnauthorized,
		},
	}

	tests[0].args.st.EXPECT().GetUserWithdrawals(ctx, "1").Return(nil, errors.New(""))
	tests[1].args.st.EXPECT().GetUserWithdrawals(ctx, "1").Return(&[]storage.Withdrawal{}, nil)
	tests[2].args.st.EXPECT().GetUserWithdrawals(ctx, "1").Return(nil, nil)
	tests[3].args.st.EXPECT().GetUserWithdrawals(ctx, "1").Return(nil, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IjEifQ.89buB88ZZu4dHpaH7b229SHhQe67gq-5Pgig2xiKm48",
			})
			w := httptest.NewRecorder()
			r := chi.NewRouter()

			st := storage.StorageRepo(tt.args.st)
			r.Post("/api/user/withdrawals", func(rw http.ResponseWriter, r *http.Request) {
				WithdrawalsGET(rw, r, &st)
			})
			r.ServeHTTP(w, tt.args.req)

			result := w.Result()
			require.Equal(t, tt.statusCode, result.StatusCode)
		})
	}

	t.Run(tests[3].name, func(t *testing.T) {
		w := httptest.NewRecorder()
		r := chi.NewRouter()

		st := storage.StorageRepo(tests[3].args.st)
		r.Post("/api/user/withdrawals", func(rw http.ResponseWriter, r *http.Request) {
			WithdrawalsGET(rw, r, &st)
		})
		r.ServeHTTP(w, tests[3].args.req)

		result := w.Result()
		require.Equal(t, tests[3].statusCode, result.StatusCode)
	})
}
