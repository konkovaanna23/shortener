package config

import (
	"flag"
	"os"
)

const (
	defaultHost     = "localhost:8080"
	defaultURLShort = "http://localhost:8080"
	defaultFilePath = "shorturl.json"
)

type Config struct {
	URLserver   string
	URLforShort string
	FilePath    string
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
	flag.Parse()

	urlServer := getEnvString("SERVER_ADDRESS", *urlServerFlag)
	urlForShort := getEnvString("BASE_URL", *urlForShortFlag)
	fileStoragePath := getEnvString("FILE_STORAGE_PATH", *fileStoragePathFlag)

	return &Config{
		URLserver:   urlServer,
		URLforShort: urlForShort,
		FilePath:    fileStoragePath,
	}
}
