package accrualmanager

import (
	"encoding/json"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/config"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/consterror"
	status "github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/status"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Manager struct {
	dateCh  chan storage.Order
	Storage *storage.Storage
	Config  *config.Config
}

func InitAccrualManager(config *config.Config, st *storage.Storage, dateCh chan storage.Order) *Manager {
	return &Manager{
		dateCh:  dateCh,
		Storage: st,
		Config:  config,
	}
}

func (m *Manager) Start() {
	client := http.Client{}
	for order := range m.dateCh {
		updatedOrder, err := sendRequestToAccrual(m.Config, order, &client)
		if err != nil {
			log.Println(err)
			m.dateCh <- order
			continue
		}
		if updatedOrder == nil {
			log.Println("The order is not registered in the accrual system")
			continue
		}
		if updatedOrder.Status != status.Processed && updatedOrder.Status != status.Invalid {
			err := m.Storage.UpdateOrder(*updatedOrder)
			if err != nil {
				m.dateCh <- order
			}
			continue
		}
		if order.Status != updatedOrder.Status {
			m.Storage.UpdateOrder(*updatedOrder)
			m.dateCh <- order
		}
	}
}

func sendRequestToAccrual(config *config.Config, order storage.Order, client *http.Client) (*storage.Order, error) {
	urlStr, err := url.JoinPath(config.GetAccrualAddr(), "/api/orders/", strconv.FormatInt(order.Number, 10))
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	var updatedOrder storage.Order
	err = retry.Retry(
		func(attempt uint) error {
			//send request
			resp, err := client.Do(r)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			switch resp.StatusCode {
			case http.StatusOK:
				{
					err := json.NewDecoder(resp.Body).Decode(&updatedOrder)
					if err != nil {
						return err
					}
				}
			case http.StatusNoContent:
				{
					log.Println("The order is not registered in the accrual system")
				}
			case http.StatusTooManyRequests:
				{
					log.Println(consterror.ErrCountRequestToAccrual)
					return consterror.ErrCountRequestToAccrual
				}
			}
			return nil
		},
		strategy.Limit(4),
		strategy.Backoff(backoff.Incremental(0, 10*time.Second)),
	)
	return &updatedOrder, err
}
