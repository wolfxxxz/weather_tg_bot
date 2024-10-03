package config

import (
	"fmt"
	"weather_bot/internal/apperrors"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Token             string `env:"TOKEN"`
	LogLevel          string `env:"LOGGER_LEVEL"`
	Host              string `env:"HOST"`
	KeyApi            string `env:"PRIMARY_KEY"`
	MongoHost         string `env:"MONGO_URL"`
	MongoPort         string `env:"MONGO_PORT"`
	UserName          string `env:"USER_NAME"`
	DBName            string `env:"DB_NAME"`
	Password          string `env:"PASSWORD"`
	TimeoutMongoQuery string `env:"TIMEOUT_MONGO_QUERY"`
}

func NewConfig() *Config {
	return &Config{}
}

func (v *Config) ParseConfig(path string, log *logrus.Logger) error {
	err := godotenv.Load(path)
	if err != nil {
		errMsg := fmt.Sprintf(" %s", err.Error())
		//return apperrors.EnvConfigParseError.AppendMessage(errMsg)
		log.Info("gotoenv could not find .env", errMsg)
	}

	if err := env.Parse(v); err != nil {
		errMsg := fmt.Sprintf("%+v\n", err)
		return apperrors.EnvConfigParseError.AppendMessage(errMsg)
	}

	log.Info("Config has been parsed, succesfully!!!")
	return nil
}
