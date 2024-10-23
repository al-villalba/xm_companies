package common

import (
	"reflect"
	"strings"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Env struct {
	AppEnv                string `env:"APP_ENV" envDefault:"dev"`
	Addr                  string `env:"ADDR" envDefault:"3000"`
	JwtKey                string `env:"JWT_KEY,required"`
	KafkaTopic            string `env:"KAFKA_TOPIC,required"`
	KafkaBootstrapServers string `env:"KAFKA_BOOTSTRAP_SERVERS,required"`
	KafkaMessageTimeout   string `env:"KAFKA_MESSAGE_TIMEOUT" envDefault:"5000"`
	KafkaRetries          string `env:"KAFKA_RETRIES" envDefault:"10"`
	KafkaRetryBackoff     string `env:"KAFKA_RETRY_BACKOFF" envDefault:"500"`
	DbHost                string `env:"DB_HOST,required"`
	DbName                string `env:"DB_NAME,required"`
	DbUser                string `env:"DB_USER,required"`
	DbPass                string `env:"DB_PASS,required"`
}

var e *Env
var parsers = map[reflect.Type]env.ParserFunc{}

func GetEnv() *Env {
	if e != nil {
		return e
	}

	e = &Env{}
	if err := env.ParseWithFuncs(e, parsers); err != nil {
		logrus.WithError(err).Panic("Unable to parse envs")
	}

	// Setup logger
	logrus.SetLevel(logrus.InfoLevel)
	if !e.IsProd() {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
	}

	return e
}

func (e *Env) IsProd() bool {
	return strings.HasPrefix(e.AppEnv, "prod")
}
