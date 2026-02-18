// Package audit - реализует аудит событий.
package audit

import (
	"encoding/json"
	"time"
)

// Event — событие аудита
type Event struct {
	Time   time.Time `json:"ts"`
	Action string    `json:"action"`
	UserID string    `json:"user_id"`
	URL    string    `json:"url"`
}

// MarshalJSON реализует интерфейс json.Marshaler.
func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		TS     int64  `json:"ts"`
		Action string `json:"action"`
		UserID string `json:"user_id"`
		URL    string `json:"url"`
	}{
		TS:     e.Time.Unix(),
		Action: e.Action,
		UserID: e.UserID,
		URL:    e.URL,
	})
}
