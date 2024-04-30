package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/config"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/consterror"
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

func RegistrationPOST(w http.ResponseWriter, r *http.Request, st *storage.Storage) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = st.AddNewUser(user.Login, user.Password)
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

func LoginPOST(w http.ResponseWriter, r *http.Request, st *storage.Storage) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	//TODO: check user password if exist
	ok, err := st.CheckUserPassword(user.Login, user.Password)
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

func OrdersPOST(w http.ResponseWriter, r *http.Request, st *storage.Storage) {
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

	_, err = st.AddNewOrder(claims.Login, orderNum)
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

func OrdersGET(w http.ResponseWriter, r *http.Request, conf *config.Config, st *storage.Storage) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	orders, err := st.GetUserOrders(claims.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := http.Client{}
	var ordersInfo []storage.Order
	for _, order := range *orders {
		orderAccrual, err := sendRequestToAccrual(conf, order, &client)
		if err != nil {
			continue
		}
		orderNumber, err := strconv.ParseInt(orderAccrual.Order, 10, 64)
		if err != nil {
			continue
		}
		updatedOrder := storage.Order{
			Number:     orderNumber,
			Status:     orderAccrual.Status,
			Accrual:    orderAccrual.Accrual,
			UploadDate: order.UploadDate,
			Sum:        order.Sum,
		}
		err = st.UpdateOrder(updatedOrder)
		if err != nil {
			return
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

	w.WriteHeader(http.StatusOK)
}

func WithdrawalsGET(w http.ResponseWriter, r *http.Request, st *storage.Storage) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	withdrawals, err := st.GetUserWithdrawals(claims.Login)
	if len(*withdrawals) == 0 {
		http.Error(w, err.Error(), http.StatusNoContent)
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
	w.WriteHeader(http.StatusOK)
}

func BalanceGET(w http.ResponseWriter, r *http.Request, st *storage.Storage) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	balance, err := st.GetUserBalance(claims.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func BalanceWithdrawPOST(w http.ResponseWriter, r *http.Request, st *storage.Storage) {
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
	//check withdrawals.Order with luna
	/*
		if not correct {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	*/

	err = st.DoRebiting(claims.Login, withdrawals.Order, withdrawals.Sum)
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
	urlStr, err := url.JoinPath(config.GetAccrualAddr(), "/api/orders/", strconv.FormatInt(order.Number, 10))
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	var updatedOrder OrderAccrual
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
