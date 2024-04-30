package gophermart

import (
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/accrualmanager"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/parser"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/router"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func InitGophermart() {
	conf := parser.ParseConfig()
	st := parser.ParseStorageInfo()
	orderCh := make(chan storage.Order, 100)
	am := accrualmanager.InitAccrualManager(conf, st, orderCh)
	go am.Start()

	r := router.InitRouter(st, orderCh)
	log.Println(conf.GetSrvAddr())
	err := http.ListenAndServe(conf.GetSrvAddr() /*conf.GetSrvAddr()*/, r)

	if err != nil {
		panic(err)
	}
}
