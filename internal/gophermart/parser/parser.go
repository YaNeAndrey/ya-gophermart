package parser

import (
	"flag"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/config"
	"github.com/YaNeAndrey/ya-gophermart/internal/gophermart/storage"
	log "github.com/sirupsen/logrus"
	"os"
)

func ParseConfig() *config.Config {
	conf := new(config.Config)

	runAddress := flag.String("a", "localhost:8080", "Server endpoint address server:port")
	accrualAddress := flag.String("r", "localhost:8081", "Accrual system address server:port")
	flag.Parse()

	runAddressEnv, isExist := os.LookupEnv("RUN_ADDRESS")
	var err error
	if !isExist {
		err = conf.SetSrvAddr(*runAddress)
	} else {
		err = conf.SetSrvAddr(runAddressEnv)
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	accrualAddressEnv, isExist := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if !isExist {
		err = conf.SetSrvAddr(*accrualAddress)
	} else {
		err = conf.SetSrvAddr(accrualAddressEnv)
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	return conf
}

func ParseStorageInfo() *storage.Storage {
	st := new(storage.Storage)
	databaseURI := flag.String("d", "", "DB connection string")

	flag.Parse()
	var err error
	databaseURIEnv, isExist := os.LookupEnv("DATABASE_URI")
	if isExist {
		err = st.SetDBConnectionString(databaseURIEnv)
	} else {
		err = st.SetDBConnectionString(*databaseURI)
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	return st
}
