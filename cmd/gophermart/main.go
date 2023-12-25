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

	config.ParseEnvFlags()
	logger.Logger.Info("Creating storage...")
	storage.Stor, err = storage.NewPgStorage()
	if err != nil {
		logger.Logger.Error(err.Error())
		panic(err)
	}
	logger.Logger.Info("Creating storage DONE")
	logger.Logger.Info("Creating router...")
	r := router.New()
	if r == nil {
		err = errors.New("failed to create main router")
		logger.Logger.Error(err.Error())
		panic(err)
	}
	logger.Logger.Info("Creating storage DONE")
	logger.Logger.Info("Start server...")

	done := make(chan bool)
	err = http.ListenAndServe(config.Options.RunAddress, r)
	if err != nil {
		logger.Logger.Error(err.Error())
		panic(err)
	}
	logger.Logger.Info("Start server DONE")
	done <- true
}
