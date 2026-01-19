package audit

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// HTTPObserver отправляет события аудита на HTTP-сервер
type HTTPObserver struct {
	client *http.Client
	url    string
}

// NewHTTPObserver создает новый HTTPObserver и возвращает его
// Параметры:
//   - url: URL HTTP-сервера, на который будут отправляться события аудита
func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{
		client: &http.Client{},
		url:    url,
	}
}

// Notify отправляет событие аудита на HTTP-сервер
func (h *HTTPObserver) Notify(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		logrus.Error(err)
		return
	}

	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(data))
	if err != nil {
		logrus.Error(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logrus.Error("Аудит-сервис вернул ошибку ", resp.StatusCode)

	}
}
