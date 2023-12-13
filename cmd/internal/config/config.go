package config

import (
	flag "github.com/spf13/pflag"
	"os"
)

var Options struct {
	RunAddress  string
	DatabaseURI string
	HMACSecret  string
}

func ParseEnvFlags() {
	var ok bool
	Options.HMACSecret, ok = os.LookupEnv("HMAC_SECRET")
	if !ok {
		flag.StringVarP(&Options.HMACSecret, "s", "s", "zdgLBLCdslbvbsVJCLDcvdhlsvlshd", "HMAC secret")
	}
	Options.DatabaseURI, ok = os.LookupEnv("DATABASE_URI")
	if !ok {
		flag.StringVarP(&Options.DatabaseURI, "d", "d", "", "Строка подключения к БД")
	}
	Options.RunAddress, ok = os.LookupEnv(`RUN_ADDRESS`)
	if !ok {
		flag.StringVarP(&Options.RunAddress, "a", "a", ":8080", "Адрес HTTP-сервера")
	}
	if !ok {
		flag.Parse()
	}
}
