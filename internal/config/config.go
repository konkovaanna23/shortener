// Package config - пакет для работы с конфигурацией.
package config

import (
	"flag"
	"os"
	"strconv"
)

const (
	defaultHost         = "localhost:8080"
	defaultURLShort     = "http://localhost:8080"
	defaultFilePath     = "shorturl.json"
	defaultDSN          = "postgres://user_main:user_main@localhost:5432/shortenerdb?sslmode=disable"
	defaultBufferSize   = 100
	defaultBatchSize    = 3
	defaultKey          = "secret"
	defaultTimeFlushDel = 2
)

// Config Конфигурация приложения.
type Config struct {
	URLserver     string
	URLforShort   string
	FilePath      string
	DSN           string
	BufferSize    int
	BatchSize     int
	Key           string
	TimeFlushDel  int
	AuditFilePath string
	AuditURL      string
}

func getEnvString(envKey, defaultValue string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt(envKey string, defaultValue int) int {
	if v := os.Getenv(envKey); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

// GetConfig возвращает конфигурацию приложения.
func GetConfig() *Config {
	urlServerFlag := flag.String("a", defaultHost, "Адрес запуска HTTP-сервера")
	urlForShortFlag := flag.String("b", defaultURLShort, "Основной URL для сокращения")
	fileStoragePathFlag := flag.String("f", "", "Путь до файла")
	dsnFlag := flag.String("d", "", "DSN для подключения к базе данных")
	bufferSizeFlag := flag.Int("u", defaultBufferSize, "Размер буфера для накопления объектов обновления")
	batchSizeFlag := flag.Int("h", defaultBatchSize, "Размер обновляемых URL для удаления")
	timeFlushDelFlag := flag.Int("t", defaultTimeFlushDel, "Период ожидания обновления удаляемых данных(в секундах)")
	keyFlag := flag.String("k", defaultKey, "Ключ для шифрования пользователя")
	auditFileFlag := flag.String("audit-file", "", "Путь до файла аудита")
	auditURLFlag := flag.String("audit-url", "", "URL для аудита")
	flag.Parse()

	urlServer := getEnvString("SERVER_ADDRESS", *urlServerFlag)
	urlForShort := getEnvString("BASE_URL", *urlForShortFlag)
	fileStoragePath := getEnvString("FILE_STORAGE_PATH", *fileStoragePathFlag)
	dsn := getEnvString("DSN", *dsnFlag)
	bufferSize := getEnvInt("BUFFER_SIZE", *bufferSizeFlag)
	batchSize := getEnvInt("BATCH_SIZE", *batchSizeFlag)
	timeFlushDel := getEnvInt("TIME_FLUSH_DELETE", *timeFlushDelFlag)
	key := getEnvString("KEY", *keyFlag)
	auditFile := getEnvString("AUDIT_FILE", *auditFileFlag)
	auditURL := getEnvString("AUDIT_URL", *auditURLFlag)

	return &Config{
		URLserver:     urlServer,
		URLforShort:   urlForShort,
		FilePath:      fileStoragePath,
		DSN:           dsn,
		BufferSize:    bufferSize,
		BatchSize:     batchSize,
		TimeFlushDel:  timeFlushDel,
		Key:           key,
		AuditFilePath: auditFile,
		AuditURL:      auditURL,
	}
}
