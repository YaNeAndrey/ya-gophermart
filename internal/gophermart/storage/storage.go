package storage

import "context"

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

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Orders (ID_order SERIAL PRIMARY KEY, status VARCHAR(10) NOT NULL, uploaded_at timestamp NOT NULL, sum float null, accrual float NULL);")
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
	row := db.QueryRowContext(ctx, "SELECT (case when (passwd = crypt($2, passwd)) then 'True' else 'False' end) as ok FROM users WHERE login = $1 ", login, password)
	var passwdOK bool

	// if error no rows - return user not found
	err = row.Scan(&passwdOK)
	if err != nil {
		return false, err
	}

	return passwdOK, nil
}

// Загрузка номера заказа
func (s *Storage) AddNewOrder(login string) error {
	return nil
}

// Получение текущего баланса пользователя
func (s *Storage) GetUserBalance(login string) (*Balance, error) {
	return nil, nil
}

// Получение списка загруженных номеров заказов
func (s *Storage) GetUserOrders(login string) (*[]Order, error) {
	return nil, nil
}

// Получение информации о выводе средств
func (s *Storage) GetUserWithdrawals(login string) (*[]Withdrawal, error) {
	return nil, nil
}

// Запрос на списание средств
func (s *Storage) DoRebiting(login string, sum float32) error {
	return nil
}
