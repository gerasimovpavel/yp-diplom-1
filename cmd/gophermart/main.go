package main

import (
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/logger"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/router"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

func main() {
	var err error
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger.Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}

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
	done := make(chan bool)
	//scheduler.Schedule(accruals.CheckAccruals, time.Second, done)

	err = http.ListenAndServe(config.Options.RunAddress, router)
	if err != nil {
		panic(err)
	}
	done <- true
}
