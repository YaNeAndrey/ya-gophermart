package gophermart

import (
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/accrualmanager"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/parser"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/router"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	"net/http"
	"time"
)

func InitGophermart() {
	conf := parser.ParseConfig()
	st := parser.ParseStorageInfo()
	orderCh := make(chan storage.Order, 100)
	orderCh <- storage.Order{
		Number:        1,
		Status:        "22",
		Accrual:       0,
		UploadDate:    time.Time{},
		Sum:           0,
		ProcessedDate: time.Time{},
	}
	orderCh <- storage.Order{
		Number:        5,
		Status:        "22",
		Accrual:       0,
		UploadDate:    time.Time{},
		Sum:           0,
		ProcessedDate: time.Time{},
	}
	am := accrualmanager.InitAccrualManager(conf, st, orderCh)
	am.Start()

	r := router.InitRouter(st, orderCh)

	err := http.ListenAndServe(conf.GetSrvAddr(), r)

	if err != nil {
		panic(err)
	}
}
