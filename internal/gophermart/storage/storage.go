package storage

import "context"

type StorageRepo interface {
	AddNewUser(ctx context.Context, login string, password string) error
	CheckUserPassword(ctx context.Context, login string, password string) (bool, error)
	AddNewOrder(ctx context.Context, login string, orderNumber string) (*Order, error)
	GetUserBalance(ctx context.Context, login string) (*Balance, error)
	GetUserOrders(ctx context.Context, login string) (*[]Order, error)
	GetUserWithdrawals(ctx context.Context, login string) (*[]Withdrawal, error)
	DoRebiting(ctx context.Context, login string, order string, sum float64) error
	GetAllNotProcessedOrders(ctx context.Context) (*[]Order, error)
	UpdateOrder(ctx context.Context, order Order) error
	UpdateBalance(ctx context.Context, order Order) error
}
