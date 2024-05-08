package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/config"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/consterror"
	status "github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/status"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Userinfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func RegistrationPOST(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = (*st).AddNewUser(ctx, user.Login, user.Password)
	if err != nil {
		if err == consterror.ErrDuplicateLogin {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = SetToken(&w, user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginPOST(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	ok, err := (*st).CheckUserPassword(ctx, user.Login, user.Password)
	if err != nil {
		if err == consterror.ErrLoginNotFound {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if ok {
		err = SetToken(&w, user.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func OrdersPOST(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	body, _ := io.ReadAll(r.Body)
	orderNum, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = goluhn.Validate(strconv.FormatInt(orderNum, 10))
	if err != nil {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}
	ctx := context.Background()
	_, err = (*st).AddNewOrder(ctx, claims.Login, strconv.FormatInt(orderNum, 10))
	if err != nil {
		switch err {
		case consterror.ErrDuplicateUserOrder:
			{
				w.WriteHeader(http.StatusOK)

			}
		case consterror.ErrDuplicateAnotherUserOrder:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func OrdersGET(w http.ResponseWriter, r *http.Request, conf *config.Config, st *storage.StorageRepo) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	orders, err := (*st).GetUserOrders(ctx, claims.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	var ordersInfo []storage.Order
	for _, order := range *orders {
		updatedOrder := order
		if order.Status != status.Processed && order.Status != status.Invalid {
			orderAccrual, err := sendRequestToAccrual(conf, order, &client)
			if err != nil {
				ordersInfo = append(ordersInfo, order)
				continue
			}
			if orderAccrual == nil {
				ordersInfo = append(ordersInfo, order)
				continue
			}
			updatedOrder = storage.Order{
				Number:     orderAccrual.Order,
				Status:     orderAccrual.Status,
				Accrual:    orderAccrual.Accrual,
				UploadDate: order.UploadDate,
				Sum:        order.Sum,
			}
		}

		if updatedOrder.Status != order.Status {
			_ = (*st).UpdateOrder(ctx, updatedOrder)
			if updatedOrder.Status == status.Processed {
				_ = (*st).UpdateBalance(ctx, updatedOrder)
			}
		}
		ordersInfo = append(ordersInfo, updatedOrder)
	}

	if len(ordersInfo) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := json.Marshal(ordersInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
}

func WithdrawalsGET(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	withdrawals, err := (*st).GetUserWithdrawals(ctx, claims.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if withdrawals == nil {
		http.Error(w, "", http.StatusNoContent)
	}
	body, err := json.Marshal(withdrawals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//w.WriteHeader(http.StatusOK)
}

func BalanceGET(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	balance, err := (*st).GetUserBalance(ctx, claims.Login)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(balance)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
}

func BalanceWithdrawPOST(w http.ResponseWriter, r *http.Request, st *storage.StorageRepo) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	var withdrawals WithdrawalsRequest
	err := json.NewDecoder(r.Body).Decode(&withdrawals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	err = (*st).DoRebiting(ctx, claims.Login, withdrawals.Order, withdrawals.Sum)
	if err != nil {
		switch err {
		case consterror.ErrOrderNotFound:
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case consterror.ErrInsufficientFunds:
			http.Error(w, err.Error(), http.StatusPaymentRequired)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ReadAuthDate(r *http.Request) (*Userinfo, error) {
	var user Userinfo
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SetToken(w *http.ResponseWriter, login string) error {
	token, err := buildJWTStringWithLogin(login)
	if err != nil {
		return err
	}
	http.SetCookie(*w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
	(*w).Header().Set("Content-Type", "application/json")
	return nil
}

func buildJWTStringWithLogin(login string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Login:            login,
		RegisteredClaims: jwt.RegisteredClaims{},
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(constants.SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func checkAccess(r *http.Request) (*Claims, bool) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, false
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(constants.SecretKey), nil
		})

	if err != nil {
		return nil, false
	}
	if token.Valid {
		return claims, true
	} else {
		return nil, false
	}
}

func sendRequestToAccrual(config *config.Config, order storage.Order, client *http.Client) (*OrderAccrual, error) {
	urlStr, err := url.JoinPath(config.GetAccrualAddr(), "/api/orders/", order.Number)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	updatedOrder := new(OrderAccrual)
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
					err := json.NewDecoder(resp.Body).Decode(updatedOrder)
					if err != nil {
						return err
					}
				}
			case http.StatusNoContent:
				{
					log.Println("The order is not registered in the accrual system")
					updatedOrder = nil
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
	return updatedOrder, err
}
