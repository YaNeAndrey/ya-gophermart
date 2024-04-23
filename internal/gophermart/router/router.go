package router

import (
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/config"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/handler"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/middleware"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	"github.com/go-chi/chi/v5"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func InitRouter(c config.Config, st storage.Storage) http.Handler {
	r := chi.NewRouter()
	r.NotFound(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	})

	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Gzip())

	r.Route("/", func(r chi.Router) {
		r.Post("/api/user/register/", func(rw http.ResponseWriter, req *http.Request) {
			handler.RegistrationPOST(rw, req)
		})
		r.Post("/api/user/login/", func(rw http.ResponseWriter, req *http.Request) {
			handler.LoginPOST(rw, req)
		})

		//after authorization
		r.Route("/api/user/", func(r chi.Router) {
			r.Use(middleware.Authorization())
			r.Route("/orders/", func(r chi.Router) {
				r.Post("/", func(rw http.ResponseWriter, req *http.Request) {
					handler.OrdersPOST(rw, req)
				})
				r.Get("/", func(rw http.ResponseWriter, req *http.Request) {
					handler.OrdersGET(rw, req)
				})
			})

			r.Get("/withdrawals/", func(rw http.ResponseWriter, req *http.Request) {
				handler.WithdrawalsGET(rw, req)
			})
			r.Route("/balance/", func(r chi.Router) {
				r.Get("/", func(rw http.ResponseWriter, req *http.Request) {
					handler.BalanceGET(rw, req)
				})
				r.Post("/withdraw/", func(rw http.ResponseWriter, req *http.Request) {
					handler.BalanceWithdrawPOST(rw, req)
				})
			})
		})
	})
	return r
}