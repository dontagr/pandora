package models

import (
	"time"
)

type (
	User struct {
		ID           int    `json:"id"`
		Login        string `json:"login"`
		PasswordHash string `json:"password"`
	}
	SecretList struct {
		List []SecretLite `json:"list"`
	}
	SecretLite struct {
		Label   string    `json:"label"`
		DT      time.Time `json:"dt"`
		Version int       `json:"version"`
	}
	Secret struct {
		Kind    string    `json:"kind"`
		Label   string    `json:"label"`
		Version int       `json:"version"`
		Data    string    `json:"data"`
		Meta    string    `json:"meta"`
		DT      time.Time `json:"dt"`
	}
)
