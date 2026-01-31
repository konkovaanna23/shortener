package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/konkovaanna23/shortener/internal/config"
	"github.com/konkovaanna23/shortener/internal/handler"
	"github.com/konkovaanna23/shortener/internal/model"
	"github.com/konkovaanna23/shortener/internal/service"
	"github.com/sirupsen/logrus"
)

// ExampleServer_StartAndListen() демонстрирует использование всех HTTP-эндпоинтов сервера:
// - POST / → создание короткого URL
// - GET /{shorturl} → редирект
// - POST /api/shorten → JSON-создание
// - POST /api/shorten/batch → пакетное создание
// - GET /api/user/urls → получение всех URL пользователя
// - DELETE /api/user/urls → асинхронное удаление
// - GET /ping → проверка БД
func ExampleServer_Start() {
	// Инициализация
	logrus.SetLevel(logrus.ErrorLevel) // уменьшаем шум

	// Создаём конфиг
	cfg := &config.Config{
		URLserver:   "http://localhost:8081",
		URLforShort: "http://localhost:8081",
	}

	// Создаём конвертер
	ctx := context.Background()
	converter := service.NewConverter(ctx, cfg.URLforShort, cfg.FilePath, nil, 100, 10, 1)

	// Создаём сервер
	server := handler.NewServer(cfg.URLserver, converter, cfg.Key, cfg.AuditFilePath, cfg.AuditURL)

	// Запускаем в фоне
	ts := httptest.NewServer(server.GetHandler())
	defer ts.Close()

	// Даем серверу время запуститься
	time.Sleep(100 * time.Millisecond)

	// === 1. POST / — создание короткого URL (текст)
	resp, err := http.Post(ts.URL+"/", "text/plain", bytes.NewBufferString("https://example.com"))
	if err != nil {
		fmt.Printf("Ошибка POST /: %v\n", err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("POST /: %s → %s (Status: %d)\n", "https://example.com", string(body), resp.StatusCode)
	_ = resp.Body.Close()

	shortURL := string(body)

	// === 2. GET /{shorturl} — редирект
	resp, err = http.Get(ts.URL + "/" + shortURL[22:]) // извлекаем /abc123
	if err != nil {
		fmt.Printf("Ошибка GET /{shorturl}: %v\n", err)
		return
	}
	fmt.Printf("GET /{shorturl}: Redirect to %s (Status: %d)\n", resp.Request.URL.String(), resp.StatusCode)
	_ = resp.Body.Close()

	// === 3. POST /api/shorten — JSON
	request := &service.URLRequest{
		URL: "https://golang.org",
	}
	jsonBody, _ := json.Marshal(request)
	resp, err = http.Post(ts.URL+"/api/shorten", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Ошибка POST /api/shorten: %v\n", err)
		return
	}
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("POST /api/shorten: %s (Status: %d)\n", string(body), resp.StatusCode)
	_ = resp.Body.Close()

	var result service.URLResponse
	_ = json.Unmarshal(body, &result)

	// === 4. POST /api/shorten/batch — пакетное создание
	batch := []*model.DescriptionURL{
		{Correlation: "id1", Original: "https://ya.ru"},
		{Correlation: "id2", Original: "https://google.com"},
	}
	jsonBody, _ = json.Marshal(batch)
	resp, err = http.Post(ts.URL+"/api/shorten/batch", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Ошибка POST /api/shorten/batch: %v\n", err)
		return
	}
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("POST /api/shorten/batch: %s (Status: %d)\n", string(body), resp.StatusCode)
	_ = resp.Body.Close()

	// === 5. GET /api/user/urls — получение всех URL пользователя
	req, _ := http.NewRequest("GET", ts.URL+"/api/user/urls", nil)
	req.Header.Set("Authorization", "some-user-id")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Ошибка GET /api/user/urls: %v\n", err)
		return
	}
	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("GET /api/user/urls: %s (Status: %d)\n", string(body), resp.StatusCode)
	} else {
		fmt.Printf("GET /api/user/urls: пусто (Status: %d)\n", resp.StatusCode)
	}
	_ = resp.Body.Close()

	// === 6. DELETE /api/user/urls — асинхронное удаление
	deleteReq := []string{result.URLShort}
	jsonBody, _ = json.Marshal(deleteReq)
	req, _ = http.NewRequest("DELETE", ts.URL+"/api/user/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "some-user-id")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Ошибка DELETE /api/user/urls: %v\n", err)
		return
	}
	fmt.Printf("DELETE /api/user/urls: принято (Status: %d)\n", resp.StatusCode)
	_ = resp.Body.Close()

	// === 7. GET /ping — проверка БД (всегда OK при работе в памяти)
	resp, err = http.Get(ts.URL + "/ping")
	if err != nil {
		fmt.Printf("Ошибка GET /ping: %v\n", err)
		return
	}
	fmt.Printf("GET /ping: %s (Status: %d)\n", http.StatusText(resp.StatusCode), resp.StatusCode)
	_ = resp.Body.Close()

}
