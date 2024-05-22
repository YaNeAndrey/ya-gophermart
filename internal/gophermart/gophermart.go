package gophermart

import (
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/parser"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/router"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	"net/http"
)

func InitGophermart() {
	conf := parser.ParseConfig()
	st := storage.StorageRepo(parser.ParseStorageInfo())

	r := router.InitRouter(&st, conf)
	err := http.ListenAndServe(conf.GetSrvAddr(), r)

	if err != nil {
		panic(err)
	}
}
