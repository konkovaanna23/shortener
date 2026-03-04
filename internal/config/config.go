// Package config - пакет для работы с конфигурацией.
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/konkovaanna23/shortener/internal/file"
)

const (
	defaultHost         = "localhost:8080"
	defaultURLShort     = "http://localhost:8080"
	defaultBufferSize   = 100
	defaultBatchSize    = 3
	defaultKey          = "secret"
	defaultTimeFlushDel = 2
)

// Config Конфигурация приложения.
type Config struct {
	URLserver     string `json:"server_address"`
	URLforShort   string `json:"base_url"`
	FilePath      string `json:"file_storage_path"`
	DSN           string `json:"dsn"`
	BufferSize    int    `json:"buffer_size"`
	BatchSize     int    `json:"batch_size"`
	Key           string `json:"key"`
	TimeFlushDel  int    `json:"time_flush_delete"`
	AuditFilePath string `json:"audit_file"`
	AuditURL      string `json:"audit_url"`
	EnableHTTPS   bool   `json:"enable_https"`
	TrustedSubnet string `json:"trusted_subnet"`
	GrpcServer    string `json:"grpc_server_address"`
}

type configPointer struct {
	URLserver     *string `json:"server_address"`
	URLforShort   *string `json:"base_url"`
	FilePath      *string `json:"file_storage_path"`
	DSN           *string `json:"dsn"`
	BufferSize    *int    `json:"buffer_size"`
	BatchSize     *int    `json:"batch_size"`
	Key           *string `json:"key"`
	TimeFlushDel  *int    `json:"time_flush_delete"`
	AuditFilePath *string `json:"audit_file"`
	AuditURL      *string `json:"audit_url"`
	EnableHTTPS   *bool   `json:"enable_https"`
	ConfigPath    *string
	TrustedSubnet *string `json:"trusted_subnet"`
	GrpcServer    *string `json:"grpc_server_address"`
}

func defaultConfig() *Config {
	return &Config{
		URLserver:     defaultHost,
		URLforShort:   defaultURLShort,
		FilePath:      "", //"shorturl.json",
		DSN:           "", //"postgres://user_main:user_main@localhost:5432/shortenerdb?sslmode=disable",
		BufferSize:    defaultBufferSize,
		BatchSize:     defaultBatchSize,
		Key:           defaultKey,
		TimeFlushDel:  defaultTimeFlushDel,
		TrustedSubnet: "", //"192.168.1.0/24",
	}
}

func lookupEnvString(key string) (string, bool) {
	v := os.Getenv(key)
	if v == "" {
		return "", false
	}
	return v, true
}

func lookupEnvInt(key string) (int, bool) {
	v := os.Getenv(key)
	if v == "" {
		return 0, false
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return i, true
}

func lookupEnvBool(key string) (bool, bool) {
	v := os.Getenv(key)
	if v == "" {
		return false, false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, false
	}
	return b, true
}

// GetConfig возвращает конфигурацию приложения.
// Приоритет: ФЛАГИ > ENV > CONFIG(JSON) > DEFAULTS.
func GetConfig() *Config {

	cfgFlag := readFlag()

	cfg := &Config{}

	configPath := ""
	flag.CommandLine.Visit(func(f *flag.Flag) {
		if f.Name == "c" {
			configPath = *cfgFlag.ConfigPath
		}
	})
	if configPath == "" {
		if v, ok := lookupEnvString("CONFIG"); ok {
			configPath = v
		}
	}

	if configPath != "" {
		fc, err := loadConfigFile(configPath)
		if err != nil {
			log.Println("Ошибка чтения файла конфига:", err.Error())
		} else {
			applyConfigFile(cfg, fc)
		}
	}

	applyEnv(cfg)

	applyFlag(cfg, cfgFlag)

	return cfg
}

func readFlag() *configPointer {
	urlServerFlag := flag.String("a", defaultHost, "Адрес запуска HTTP-сервера")
	urlForShortFlag := flag.String("b", defaultURLShort, "Основной URL для сокращения")
	fileStoragePathFlag := flag.String("f", "", "Путь до файла")
	dsnFlag := flag.String("d", "", "DSN для подключения к базе данных")
	bufferSizeFlag := flag.Int("u", defaultBufferSize, "Размер буфера для накопления объектов обновления")
	batchSizeFlag := flag.Int("h", defaultBatchSize, "Размер обновляемых URL для удаления")
	timeFlushDelFlag := flag.Int("td", defaultTimeFlushDel, "Период ожидания обновления удаляемых данных(в секундах)")
	keyFlag := flag.String("k", defaultKey, "Ключ для шифрования пользователя")
	auditFileFlag := flag.String("audit-file", "", "Путь до файла аудита")
	auditURLFlag := flag.String("audit-url", "", "URL для аудита")
	enableHTTPSFlag := flag.Bool("s", false, "Включить HTTPS")
	configJSONFlag := flag.String("c", "", "Файл конфигурации")
	trustedSubnetFlag := flag.String("t", "", "Cтроковое представление бесклассовой адресации (CIDR)")
	grpcServerFlag := flag.String("g", "", "Адрес gRPC-сервера")

	flag.Parse()

	return &configPointer{
		URLserver:     urlServerFlag,
		URLforShort:   urlForShortFlag,
		FilePath:      fileStoragePathFlag,
		DSN:           dsnFlag,
		BufferSize:    bufferSizeFlag,
		BatchSize:     batchSizeFlag,
		TimeFlushDel:  timeFlushDelFlag,
		Key:           keyFlag,
		AuditFilePath: auditFileFlag,
		AuditURL:      auditURLFlag,
		EnableHTTPS:   enableHTTPSFlag,
		ConfigPath:    configJSONFlag,
		TrustedSubnet: trustedSubnetFlag,
		GrpcServer:    grpcServerFlag,
	}
}

func applyEnv(cfg *Config) {
	if v, ok := lookupEnvString("SERVER_ADDRESS"); ok {
		cfg.URLserver = v
	}
	if v, ok := lookupEnvString("BASE_URL"); ok {
		cfg.URLforShort = v
	}
	if v, ok := lookupEnvString("FILE_STORAGE_PATH"); ok {
		cfg.FilePath = v
	}
	if v, ok := lookupEnvString("DSN"); ok {
		cfg.DSN = v
	}
	if v, ok := lookupEnvInt("BUFFER_SIZE"); ok {
		cfg.BufferSize = v
	}
	if v, ok := lookupEnvInt("BATCH_SIZE"); ok {
		cfg.BatchSize = v
	}
	if v, ok := lookupEnvInt("TIME_FLUSH_DELETE"); ok {
		cfg.TimeFlushDel = v
	}
	if v, ok := lookupEnvString("KEY"); ok {
		cfg.Key = v
	}
	if v, ok := lookupEnvString("AUDIT_FILE"); ok {
		cfg.AuditFilePath = v
	}
	if v, ok := lookupEnvString("AUDIT_URL"); ok {
		cfg.AuditURL = v
	}
	if v, ok := lookupEnvBool("ENABLE_HTTPS"); ok {
		cfg.EnableHTTPS = v
	}
	if v, ok := lookupEnvString("TRUSTED_SUBNET"); ok {
		cfg.TrustedSubnet = v
	}
	if v, ok := lookupEnvString("GRPC_SERVER_ADDRESS"); ok {
		cfg.GrpcServer = v
	}
}

func applyFlag(dst *Config, src *configPointer) {
	flag.CommandLine.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "a":
			dst.URLserver = *src.URLserver
		case "b":
			dst.URLforShort = *src.URLforShort
		case "f":
			dst.FilePath = *src.FilePath
		case "d":
			dst.DSN = *src.DSN
		case "u":
			dst.BufferSize = *src.BufferSize
		case "h":
			dst.BatchSize = *src.BatchSize
		case "td":
			dst.TimeFlushDel = *src.TimeFlushDel
		case "k":
			dst.Key = *src.Key
		case "audit-file":
			dst.AuditFilePath = *src.AuditFilePath
		case "audit-url":
			dst.AuditURL = *src.AuditURL
		case "s":
			dst.EnableHTTPS = *src.EnableHTTPS
		case "t":
			dst.TrustedSubnet = *src.TrustedSubnet
		case "g":
			dst.GrpcServer = *src.GrpcServer
		}
	})
}

func loadConfigFile(jsonFile string) (*configPointer, error) {
	data, err := file.ReadFromFile(jsonFile)
	if err != nil {
		return nil, err
	}

	var fc configPointer
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, err
	}
	return &fc, nil
}

func applyConfigFile(dst *Config, src *configPointer) {
	if src.URLserver != nil {
		dst.URLserver = *src.URLserver
	}
	if src.URLforShort != nil {
		dst.URLforShort = *src.URLforShort
	}
	if src.FilePath != nil {
		dst.FilePath = *src.FilePath
	}
	if src.DSN != nil {
		dst.DSN = *src.DSN
	}
	if src.BufferSize != nil {
		dst.BufferSize = *src.BufferSize
	}
	if src.BatchSize != nil {
		dst.BatchSize = *src.BatchSize
	}
	if src.Key != nil {
		dst.Key = *src.Key
	}
	if src.TimeFlushDel != nil {
		dst.TimeFlushDel = *src.TimeFlushDel
	}
	if src.AuditFilePath != nil {
		dst.AuditFilePath = *src.AuditFilePath
	}
	if src.AuditURL != nil {
		dst.AuditURL = *src.AuditURL
	}
	if src.EnableHTTPS != nil {
		dst.EnableHTTPS = *src.EnableHTTPS
	}
	if src.TrustedSubnet != nil {
		dst.TrustedSubnet = *src.TrustedSubnet
	}
	if src.GrpcServer != nil {
		dst.GrpcServer = *src.GrpcServer
	}
}
