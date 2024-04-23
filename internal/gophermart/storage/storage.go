package storage

import "github.com/YaNeAndrey/ya-gophermart/internal/gophermart/utils"

type Storage struct {
	dbConnectionString string
}

func (c *Storage) SetDBConnectionString(dbConnectionString string) error {
	db, err := utils.TryToOpenDBConnection(dbConnectionString)
	if err != nil {
		return err
	}
	db.Close()
	c.dbConnectionString = dbConnectionString
	return nil
}

// Регистрация пользователя
func (c *Storage) AddNewUser(login string) error {
	return nil
}

// Аутентификация пользователя
func (c *Storage) CheckUserPassword(login string, password string) error {
	return nil
}

// Загрузка номера заказа
func (c *Storage) AddNewOrder(login string) error {
	return nil
}

// Получение текущего баланса пользователя
func (c *Storage) GetUserBalance(login string) (*Balance, error) {
	return nil, nil
}

// Получение списка загруженных номеров заказов
func (c *Storage) GetUserOrders(login string) (*[]Order, error) {
	return nil, nil
}

// Получение информации о выводе средств
func (c *Storage) GetUserWithdrawals(login string) (*[]Withdrawal, error) {
	return nil, nil
}

// Запрос на списание средств
func (c *Storage) DoRebiting(login string, sum float32) error {
	return nil
}
