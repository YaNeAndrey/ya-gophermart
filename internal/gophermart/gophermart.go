package gophermart

import (
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/parser"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/router"
	"net/http"
)

func InitGophermart() {
	conf := parser.ParseConfig()
	st := parser.ParseStorageInfo()

	r := router.InitRouter(st, conf)
	err := http.ListenAndServe(conf.GetSrvAddr(), r)

	if err != nil {
		panic(err)
	}
}
