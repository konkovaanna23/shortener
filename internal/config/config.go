package config

import (
	"flag"
)

const (
	defaultHost     = "localhost:8080"
	defaultURLShort = "http://localhost:8080"
)

type Config struct {
	URLserver   string
	URLforShort string
}

func GetConfig() *Config {
	URLserver := flag.String("a", defaultHost, "Адрес запуска HTTP-сервера")
	URLforShort := flag.String("b", defaultURLShort, "Основной URL для сокращения")
	flag.Parse()
	return &Config{
		URLserver:   *URLserver,
		URLforShort: *URLforShort,
	}
}
