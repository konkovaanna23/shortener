package config

import (
	"flag"
	"os"
)

const (
	defaultHost     = "localhost:8080"
	defaultURLShort = "http://localhost:8080"
	defaultFilePath = "shorturl.json"
	defaultDSN      = "postgres://user_main:user_main@localhost:5432/shortenerdb?sslmode=disable"
)

type Config struct {
	URLserver   string
	URLforShort string
	FilePath    string
	DSN         string
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
	fileStoragePathFlag := flag.String("f", defaultFilePath, "Путь до файла")
	dsnFlag := flag.String("d", defaultDSN, "DSN для подключения к базе данных")
	flag.Parse()

	urlServer := getEnvString("SERVER_ADDRESS", *urlServerFlag)
	urlForShort := getEnvString("BASE_URL", *urlForShortFlag)
	fileStoragePath := getEnvString("FILE_STORAGE_PATH", *fileStoragePathFlag)
	dsn := getEnvString("DSN", *dsnFlag)

	return &Config{
		URLserver:   urlServer,
		URLforShort: urlForShort,
		FilePath:    fileStoragePath,
		DSN:         dsn,
	}
}
