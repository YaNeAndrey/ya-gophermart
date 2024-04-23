package gophermart

import (
	"fmt"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/parser"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/router"
	"net/http"
)

func InitGophermart() {
	conf := parser.ParseConfig()
	st := parser.ParseStorageInfo()

	r := router.InitRouter(*conf, *st)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", "localhost", 8080), r)

	if err != nil {
		panic(err)
	}
}
