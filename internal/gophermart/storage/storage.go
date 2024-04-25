package storage

import (
	"context"
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
	db.Close()
	var st Storage
	st.dbConnectionString = dbConnectionString

	myContext := context.TODO()

	_, err = db.ExecContext(myContext, "CREATE TABLE Users( login VARCHAR(30) PRIMARY KEY, passwd TEXT NOT NULL, current_balance float NOT NULL, withdrawn_balance float NOT NULL);")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Orders (ID_order SERIAL PRIMARY KEY, status VARCHAR(10) NOT NULL, uploaded_at timestamp NOT NULL,processed_at timestamp, sum float null, accrual float NULL);")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(myContext, "CREATE TABLE Users_Orders( ID_user_order SERIAL PRIMARY KEY, login VARCHAR(30) REFERENCES Users (login), ID_order INTEGER REFERENCES Orders (ID_order));")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(myContext, "CREATE EXTENSION pgcrypto;")
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
	row := db.QueryRowContext(ctx, "INSERT INTO Users(login,passwd,current_balance,withdrawn_balance) values($1, $2, 0, 0) RETURNING login;", login, password)
	var bufLogin string

	err = row.Scan(&bufLogin)
	if err != nil {
		return err
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
		return false, err
	}

	return passwdOK, nil
}

// Загрузка номера заказа
func (s *Storage) AddNewOrder(login string, orderNumber string) error {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return err
	}

	ctx := context.Background()
	res, err := db.ExecContext(ctx, "INSERT INTO orders (id_order,status,uploaded_at,sum,accrual) VALUES ($1,'NEW',$2,0,0)", orderNumber, time.Now())
	if err != nil {
		if true { //err message == duplicate key => Check who added order {
			row := db.QueryRowContext(ctx, "select (case when (login = $1) then 'True' else 'False' end) from users_orders where id_order = $2", login, orderNumber)
			var isCurrentUser bool
			err = row.Scan(&isCurrentUser)
			if err != nil {
				return err
			}
			if isCurrentUser {
				//return error : current user created order earlier
			} else {
				//return error : another user created order earlier
			}
		} else {
			return err
		}
	}

	rows, err := res.RowsAffected()
	if rows == 1 {
		_, err = db.ExecContext(ctx, "INSERT INTO users_orders (id_order,login) values ($1,$2)", orderNumber, login)
		if err != nil {
			return err
		}
	}
	return nil
}

// Получение текущего баланса пользователя
func (s *Storage) GetUserBalance(login string) (*Balance, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	row := db.QueryRowContext(ctx, "select current_balance,withdrawn_balance from users where login = $1", login)
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
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order); err != nil {
			continue
		}
		orders = append(orders, order)
	}
	if len(orders) > 0 {
		return &orders, nil
	} else {
		return nil, nil
	}
}

// Получение информации о выводе средств
func (s *Storage) GetUserWithdrawals(login string) (*[]Withdrawal, error) {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "select orders.id_order,orders.sum,orders.processed_at from orders join users_orders on orders.id_order = users_orders.id_order where login = $1", login)
	defer rows.Close()

	var withdrawals []Withdrawal
	for rows.Next() {
		var withdrawal Withdrawal
		if err := rows.Scan(&withdrawal); err != nil {
			continue
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	if len(withdrawals) > 0 {
		return &withdrawals, nil
	} else {
		return nil, nil
	}
}

// Запрос на списание средств
func (s *Storage) DoRebiting(login string, sum float32) error {
	db, err := TryToOpenDBConnection(s.dbConnectionString)
	if err != nil {
		return err
	}

	ctx := context.Background()
	row := db.QueryRowContext(ctx, "select current_balance from users where login = $1", login)
	var currentBalance float32
	err = row.Scan(&currentBalance)
	if err != nil {
		return err
	}

	if currentBalance <= sum {
		return err //error wit message - не досаточно средства для списания
	} else {
		_, err = db.ExecContext(ctx, "UPDATE Users set current_balance = current_balance-$1,withdrawn_balance =withdrawn_balance+$1 where login = $2", sum, login)
		if err != nil {
			return err
		}
		//TODO. Update Orders

	}
	return nil
}
