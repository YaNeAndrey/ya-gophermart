package handler

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

type Userinfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

// TODO
const SECRET_KEY = "very_secret"

func RegistrationPOST(w http.ResponseWriter, r *http.Request) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: add user to DB if not exist

	err = SetToken(&w, user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	user, err := ReadAuthDate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//TODO: check user password if exist

	err = SetToken(&w, user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func OrdersPOST(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("OrdersPOST"))
}

func OrdersGET(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OrdersGET"))
}

func WithdrawalsGET(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WithdrawalsGET"))
}

func BalanceGET(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("BalanceGET"))
}

func BalanceWithdrawPOST(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("BalanceWithdrawPOST"))
}

func ReadAuthDate(r *http.Request) (*Userinfo, error) {
	var user Userinfo
	err := json.NewDecoder(strings.NewReader(r.Header.Get("Authorization"))).Decode(&user)
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
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
