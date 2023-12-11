package main

import (
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/router"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/pkg/errors"
	"net/http"
)

func main() {
	var err error
	//Парсим переменные и аргументы команднй строки
	config.ParseEnvFlags()
	// создаем Storage
	storage.Stor, err = storage.NewPgStorage()
	if err != nil {
		panic(err)
	}
	// запускаем сервер
	router := router.MainRouter()
	if router == nil {
		panic(errors.New("failed to create main router"))
	}
	err = http.ListenAndServe(config.Options.RunAddress, router)
	if err != nil {
		panic(err)
	}

}
