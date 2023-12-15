package config

import (
	flag "github.com/spf13/pflag"
	"os"
)

const HMACSecret string = "zdgLBLCdslbvbsVJCLDcvdhlsvlshd"

var Options struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func ParseEnvFlags() {
	var ok bool
	Options.AccrualSystemAddress, ok = os.LookupEnv(`ACCRUAL_SYSTEM_ADDRESS`)
	if !ok {
		flag.StringVarP(&Options.AccrualSystemAddress, "r", "r", ":8080", "Адрес HTTP-сервера системы начислений")
	}
	Options.DatabaseURI, ok = os.LookupEnv("DATABASE_URI")
	if !ok {
		flag.StringVarP(&Options.DatabaseURI, "d", "d", "", "Строка подключения к БД")
	}
	Options.RunAddress, ok = os.LookupEnv(`RUN_ADDRESS`)
	if !ok {
		flag.StringVarP(&Options.RunAddress, "a", "a", ":8081", "Адрес HTTP-сервера")
	}
	if !ok {
		flag.Parse()
	}
}
