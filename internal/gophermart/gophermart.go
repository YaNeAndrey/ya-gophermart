package gophermart

import (
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/accrualmanager"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/parser"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/router"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	"net/http"
)

func InitGophermart() {
	conf := parser.ParseConfig()
	st := parser.ParseStorageInfo()
	orderCh := make(chan storage.Order, 100)
	am := accrualmanager.InitAccrualManager(conf, st, orderCh)
	am.Start()

	r := router.InitRouter(st, orderCh)

	err := http.ListenAndServe(conf.GetSrvAddr(), r)

	if err != nil {
		panic(err)
	}
}
