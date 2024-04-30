package handler

import (
	"encoding/json"
	"fmt"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/constants/consterror"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"strconv"
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

func OrdersPOST(w http.ResponseWriter, r *http.Request, st *storage.Storage, orderCh chan<- storage.Order) {
	claims, ok := checkAccess(r)
	if !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	body, _ := io.ReadAll(r.Body)

	//if body str not number
	/*user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	*/
	orderNum, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order, err := st.AddNewOrder(claims.Login, orderNum)
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
	orderCh <- *order
	w.WriteHeader(http.StatusAccepted)
}

func OrdersGET(w http.ResponseWriter, r *http.Request, st *storage.Storage) {
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
	if len(*orders) == 0 {
		http.Error(w, "", http.StatusNoContent)
		return
	}

	body, err := json.Marshal(orders)
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
