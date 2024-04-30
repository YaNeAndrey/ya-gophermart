package storage

import (
	"context"
	"errors"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/consterror"
	status "github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/status"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type Storage struct {
	dbConnectionString string
}

func InitStorage(dbConnectionString string) (*Storage, error) {
	db, err := TryToOpenDBConnection(dbConnectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var st Storage
	st.dbConnectionString = dbConnectionString

	myContext := context.TODO()

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users( login VARCHAR(30) PRIMARY KEY, passwd TEXT NOT NULL, current_balance float NOT NULL, withdrawn_balance float NOT NULL);")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Orders (ID_order bigserial PRIMARY KEY, status VARCHAR(10) NOT NULL, uploaded_at timestamp NOT NULL,processed_at timestamp, sum float null, accrual float NULL);")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_Orders( ID_user_order SERIAL PRIMARY KEY, login VARCHAR(30) REFERENCES Users (login), ID_order bigint REFERENCES Orders (ID_order));")
	if err != nil {
		return nil, err
	}

	_, _ = db.ExecContext(myContext, "CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	if err != nil {
		return nil, err
	}

	return &st, nil
}

// Регистрация пользователя
func (s *Storage) AddNewUser(login string, password string) error {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return err
	}

	ctx := context.Background()
	row := db.QueryRowContext(ctx, "INSERT INTO Users(login,passwd,current_balance,withdrawn_balance) values($1, crypt($2, gen_salt('bf')), 0, 0) RETURNING login;", login, password)

	err = row.Err()
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return consterror.DuplicateLogin
			}
		} else {
			return err
		}
	}
	return nil
}

// Аутентификация пользователя
func (s *Storage) CheckUserPassword(login string, password string) (bool, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	row := db.QueryRowContext(ctx, "SELECT (case when (passwd = crypt($2, passwd)) then 'True' else 'False' end) as ok FROM users WHERE login = $1", login, password)

	var passwdOK bool
	// if error no rows - return user not found
	err = row.Scan(&passwdOK)
	if err != nil {
		return false, consterror.LoginNotFound
	}

	return passwdOK, nil
}

// Загрузка номера заказа
func (s *Storage) AddNewOrder(login string, orderNumber int64) (*Order, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	res, err := db.ExecContext(ctx, "INSERT INTO orders (id_order,status,uploaded_at,sum,accrual) VALUES ($1,$2,$3,0,0)", orderNumber, status.NEW, time.Now())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				row := db.QueryRowContext(ctx, "select (case when (login = $1) then 'True' else 'False' end) from users_orders where id_order = $2", login, orderNumber)
				if row.Err() != nil {
					return nil, err
				}
				var isCurrentUser bool
				err = row.Scan(&isCurrentUser)
				if isCurrentUser {
					return nil, consterror.DuplicateUserOrder
				} else {
					return nil, consterror.DuplicateAnotherUserOrder
				}
			}
		} else {
			return nil, err
		}
	}
	rows, err := res.RowsAffected()
	if rows == 1 {
		_, err = db.ExecContext(ctx, "INSERT INTO users_orders (id_order,login) values ($1,$2)", orderNumber, login)
		if err != nil {
			return nil, err
		}
	}
	return &Order{
		Number:        orderNumber,
		Status:        status.NEW,
		Accrual:       0,
		UploadDate:    time.Now(),
		Sum:           0,
		ProcessedDate: time.Time{},
	}, nil
}

// Получение текущего баланса пользователя
func (s *Storage) GetUserBalance(login string) (*Balance, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	row := db.QueryRowContext(ctx, "select current_balance,withdrawn_balance from users where login = $1", login)
	if row.Err() != nil {
		return nil, err
	}
	var balance Balance
	// if error no rows - return user not found
	err = row.Scan(&balance)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// Получение списка загруженных номеров заказов
func (s *Storage) GetUserOrders(login string) (*[]Order, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "select id_order,status,uploaded_at,accrual from orders join users_orders on orders.id_order = users_orders.id_order where login = $1", login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.Number, &order.Status, &order.UploadDate, &order.Accrual); err != nil {
			continue
		}
		orders = append(orders, order)
	}
	return &orders, nil

}

// Получение информации о выводе средств
func (s *Storage) GetUserWithdrawals(login string) (*[]Withdrawal, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "select orders.id_order,orders.sum,orders.processed_at from orders join users_orders on orders.id_order = users_orders.id_order where login = $1 and sum > 0", login)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	var withdrawals []Withdrawal
	for rows.Next() {
		var withdrawal Withdrawal
		if err := rows.Scan(&withdrawal.OrderNumber, &withdrawal.Sum, withdrawal.ProcessedDate); err != nil {
			continue
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	return &withdrawals, nil
}

// Запрос на списание средств
func (s *Storage) DoRebiting(login string, order int64, sum float32) error {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return err
	}
	ctx := context.Background()
	row := db.QueryRowContext(ctx, "select current_balance from users where login = $1", login)
	if row.Err() != nil {
		return err
	}
	var currentBalance float32
	err = row.Scan(&currentBalance)
	if err != nil {
		return err
	}
	if currentBalance <= sum {
		return consterror.InsufficientFunds
	}
	res, err := db.ExecContext(ctx, "update orders set sum = sum+$1 where id_order = $2", sum, order)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return consterror.OrderNotFound
	}

	_, err = db.ExecContext(ctx, "UPDATE Users set current_balance = current_balance-$1,withdrawn_balance = withdrawn_balance+$1 where login = $2", sum, login)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetAllNotProcessedOrders() (*[]Order, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "select id_order,status from orders where orders.status != 'PROCESSED' and orders.status != 'INVALID'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.Number, &order.Status); err != nil {
			continue
		}
		orders = append(orders, order)
	}
	return &orders, nil
}

func (s *Storage) UpdateOrder(order Order) error {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return err
	}
	ctx := context.Background()
	if order.Status == status.PROCESSED {
		_, err := db.ExecContext(ctx, "update orders set status = $1, processed_at = $2 where id_order = $3", order.Status, order.ProcessedDate, order.Number)
		if err != nil {
			return err
		}
	} else {
		_, err := db.ExecContext(ctx, "update orders set status = $1 where id_order = $2", order.Status, order.Number)
		if err != nil {
			return err
		}
	}
	return nil
}
