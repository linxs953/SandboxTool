package config

import (
	"os"
	"runtimeops/errs"
	"runtimeops/utils"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type EnvirConfig struct {
	TOKENBYUNIONIDAPI string `env:"TOKENBYUNIONIDAPI" envDefault:""`
	TOKENBYOPENAPI    string `env:"TOKENBYOPENAPI" envDefault:""`
	SECRETKEY         string `env:"SECRETKEY" envDefault:"secret"`
	BACKENDCICD       string `env:"BACKENDCICD" envDefault:""`
	PROD              string `env:"PROD" envDefault:""`
	TEST              string `env:"TEST" envDefault:""`
	LOCAL             string `env:"LOCAL" envDefault:""`
	RECYCLEAPI        string `env:"RECYCLEAPI" envDefault:""`
	ENTRTV2API        string `env:"ENTRYV2SB" envDefault:""`
	SANDBOXPOOL       string `env:"SANDBOXPOOL" envDefault:""`
	EntryV1           string `env:"EntryV1SB" envDefault:""`
}

var config = EnvirConfig{}

// Load config
// 1. Firstly,load current dir env file name `./.env`,if env file not exist , do 2
// 2. Load EnvirConfig in  config/config.go,then  extract key and value for setting env.
//    when setting env, if key has value(not none) do nothing
func Init() error {
	var err error
	if _, err := os.Stat("./.env"); !os.IsNotExist(err) {
		err = godotenv.Load()
		if err != nil {
			return err
		}
		return nil
	}
	if err = env.Parse(&config); err != nil {
		log.Error().Err(err).Msg("Parse env struct error")
		return err
	}
	fieldList := utils.GetTagName(config)
	value := utils.GetTagValue(config)
	err = SetEnv(fieldList, value)
	if err != nil {
		return err
	}
	return nil
}

func SetEnv(key []string, value []string) error {
	var lengthError = errs.LengthNotEqual
	var err error
	var envValue string
	if len(key) != len(value) {
		return lengthError
	}
	for i := 0; i < len(key); i++ {
		v := os.Getenv(key[i])
		if v == "" {
			envValue = value[i]
		} else {
			envValue = v
		}
		err = os.Setenv(key[i], envValue)
		if err != nil {
			log.Error().Err(err).Msg("Set os env error")
			return err
		}
	}
	return nil
}
