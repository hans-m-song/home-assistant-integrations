package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	BuildCommit = "unknown"
	BuildTime   = "unknown"

	LogLevel       string
	LogFormat      string
	ConfigFilename string
	ListenPort     int
	MQTTURI        string

	ProviderConfig = map[string]map[string]any{}
)

func Configure() {
	flag.StringVar(&LogLevel, "log-level", "info", "log level")
	flag.StringVar(&LogFormat, "log-format", "text", "log format, one of 'json', 'text'")
	flag.StringVar(&ConfigFilename, "config", "", "config file")
	flag.IntVar(&ListenPort, "listen-port", 8080, "listen port")
	flag.StringVar(&MQTTURI, "mqtt-uri", "", "MQTT URI")
	flag.Parse()

	raw, err := os.ReadFile(ConfigFilename)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	if err := yaml.Unmarshal(raw, &ProviderConfig); err != nil {
		panic(fmt.Errorf("failed to unmarshal config file: %s", err))
	}

	logLevel, err := zerolog.ParseLevel(LogLevel)
	if err != nil {
		panic(fmt.Errorf("failed to parse log level: %s", err))
	}

	zerolog.SetGlobalLevel(logLevel)

	switch LogFormat {
	case "text":
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	case "json":
		// use default logger
	default:
		panic(fmt.Errorf("unknown log format: %s", LogFormat))
	}

	log.Logger = log.Logger.With().Timestamp().Stack().Caller().Logger()
}

func Values() map[string]any {
	providers := []string{}
	for name := range ProviderConfig {
		providers = append(providers, name)
	}

	return map[string]any{
		"build_commit":    BuildCommit,
		"build_time":      BuildTime,
		"log_level":       LogLevel,
		"log_format":      LogFormat,
		"config_filename": ConfigFilename,
		"mqtt_uri":        MQTTURI,
		"providers":       providers,
	}
}
