package model

type Stats struct {
	URLs  int `json:"urls" db:"urls"`
	Users int `json:"users" db:"users"`
}
