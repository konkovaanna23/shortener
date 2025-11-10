package config

import (
	"flag"
	"os"
)

const (
	defaultHost     = "localhost:8080"
	defaultURLShort = "http://localhost:8080"
)

type Config struct {
	URLserver   string
	URLforShort string
}

func getEnvString(envKey, defaultValue string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return defaultValue
}

func GetConfig() *Config {
	urlServerFlag := flag.String("a", defaultHost, "Адрес запуска HTTP-сервера")
	urlForShortFlag := flag.String("b", defaultURLShort, "Основной URL для сокращения")
	flag.Parse()

	urlServer := getEnvString("SERVER_ADDRESS", *urlServerFlag)
	urlForShort := getEnvString("BASE_URL", *urlForShortFlag)

	return &Config{
		URLserver:   urlServer,
		URLforShort: urlForShort,
	}
}
