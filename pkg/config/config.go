package config

import (
	"fmt"
	"os"

	"github.com/codingconcepts/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	BuildCommit  = "unknown"
	BuildTime    = "unknown"
	BuildVersion = "unknown"

	Values = struct {
		LogLevel   string `env:"LOG_LEVEL" default:"info"`
		LogFormat  string `env:"LOG_FORMAT" default:"json"`
		ListenPort int    `env:"LISTEN_PORT" default:"8080"`
		MQTTURI    string `env:"MQTT_URI" required:"true"`
	}{}
)

func Configure() {
	if err := godotenv.Load(); err != nil {
		log.Debug().Err(fmt.Errorf("failed to load .env file: %s", err)).Send()
	}

	if err := env.Set(&Values); err != nil {
		panic(fmt.Errorf("failed to set env values: %s", err))
	}

	logLevel, err := zerolog.ParseLevel(Values.LogLevel)
	if err != nil {
		panic(fmt.Errorf("failed to parse log level: %s", err))
	}

	zerolog.SetGlobalLevel(logLevel)

	switch Values.LogFormat {
	case "text":
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
			With().Timestamp().Stack().Caller().Logger()
	case "json":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Logger.
			With().Stack().Caller().Logger()
	default:
		panic(fmt.Errorf("unknown log format: %s", Values.LogFormat))
	}
}
