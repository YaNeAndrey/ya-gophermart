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
	accrualAddress := flag.String("r", "http://localhost:8081", "Accrual system address server:port")
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
		err = conf.SetAccrualAddr(*accrualAddress)
	} else {
		err = conf.SetAccrualAddr(accrualAddressEnv)
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	return conf
}

func ParseStorageInfo() *storage.Storage {
	databaseURI := flag.String("d", "", "DB connection string")

	flag.Parse()
	var err error
	databaseURIEnv, isExist := os.LookupEnv("DATABASE_URI")
	st := new(storage.Storage)
	if isExist {
		st, err = storage.InitStorage(databaseURIEnv)
	} else {
		st, err = storage.InitStorage(*databaseURI)
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	return st
}
