package main

import (
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/accruals"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/logger"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/router"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/scheduler"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

func main() {
	var err error
	c := zap.NewDevelopmentEncoderConfig()
	c.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger.Logger = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(c),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
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
	scheduler.Schedule(accruals.CheckAccruals, 5*time.Second, done)

	err = http.ListenAndServe(config.Options.RunAddress, router)
	if err != nil {
		panic(err)
	}
	done <- true
}
