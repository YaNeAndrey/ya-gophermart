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

	_, err = db.ExecContext(myContext, "CREATE EXTENSION IF NOT EXISTS pgcrypto;")
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
				return consterror.ErrDuplicateLogin
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
		return false, consterror.ErrLoginNotFound
	}

	return passwdOK, nil
}

// Загрузка номера заказа
func (s *Storage) AddNewOrder(login string, orderNumber string) (*Order, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	res, err := db.ExecContext(ctx, "INSERT INTO orders (id_order,status,uploaded_at,sum,accrual) VALUES ($1,$2,$3,0,0)", orderNumber, status.New, time.Now())
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
				if err != nil {
					return nil, err
				}
				if isCurrentUser {
					return nil, consterror.ErrDuplicateUserOrder
				} else {
					return nil, consterror.ErrDuplicateAnotherUserOrder
				}
			}
		} else {
			return nil, err
		}
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 1 {
		_, err = db.ExecContext(ctx, "INSERT INTO users_orders (id_order,login) values ($1,$2)", orderNumber, login)
		if err != nil {
			return nil, err
		}
	}
	return &Order{
		Number:     orderNumber,
		Status:     status.New,
		Accrual:    0,
		UploadDate: time.Now(),
		Sum:        0,
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
	err = row.Scan(&balance.Current, &balance.Withdrawn)
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
	rows, err := db.QueryContext(ctx, "select orders.id_order,status,uploaded_at,accrual from orders join users_orders on orders.id_order = users_orders.id_order where login = $1", login)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
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
	/*
		rows, err := db.QueryContext(ctx, "select orders.id_order,orders.sum,orders.uploaded_at from orders join users_orders on orders.id_order = users_orders.id_order where login = $1", login)
		if err != nil {
			return nil, err
		}*/

	rows, err := db.QueryContext(ctx, "select orders.id_order,orders.sum,orders.uploaded_at from orders join users_orders on orders.id_order = users_orders.id_order where login = $1 and sum > 0", login)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()

	var withdrawals []Withdrawal
	for rows.Next() {
		var withdrawal Withdrawal
		if err := rows.Scan(&withdrawal.OrderNumber, &withdrawal.Sum, &withdrawal.ProcessedDate); err != nil {
			continue
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	if len(withdrawals) == 0 {
		return nil, nil
	} else {
		return &withdrawals, nil
	}
}

// Запрос на списание средств
func (s *Storage) DoRebiting(login string, order string, sum float64) error {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return err
	}
	ctx := context.Background()
	row := db.QueryRowContext(ctx, "select current_balance from users where login = $1", login)
	if row.Err() != nil {
		return err
	}
	var currentBalance float64
	err = row.Scan(&currentBalance)
	if err != nil {
		return err
	}
	if currentBalance <= sum {
		return consterror.ErrInsufficientFunds
	}
	res, err := db.ExecContext(ctx, "insert into orders (id_order,status,uploaded_at ,sum,accrual) values ($1,'NEW',$2,$3,0) ON CONFLICT (id_order) DO UPDATE SET sum = orders.sum + $3", order, time.Now(), sum)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return consterror.ErrOrderNotFound
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
	if rows.Err() != nil {
		return nil, rows.Err()
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

	_, err = db.ExecContext(ctx, "update orders set status = $1 where id_order = $2", order.Status, order.Number)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateBalance(order Order) error {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = db.ExecContext(ctx, "update Users set current_balance = current_balance+$1 where login = (select DISTINCT users_orders.login from users_orders where users_orders.id_order = $2)", order.Accrual, order.Number)
	if err != nil {
		return err
	}
	return nil
}
