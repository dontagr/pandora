package models

import (
	"time"
)

type (
	User struct {
		Login        string `json:"login"`
		PasswordHash string `json:"password"`
		ID           int    `json:"id"`
	}
	SecretList struct {
		List []SecretLite `json:"list"`
	}
	SecretLite struct {
		DT      time.Time `json:"dt"`
		Label   string    `json:"label"`
		Version int       `json:"version"`
	}
	Secret struct {
		DT      time.Time `json:"dt"`
		Kind    string    `json:"kind"`
		Label   string    `json:"label"`
		Data    string    `json:"data"`
		Meta    string    `json:"meta"`
		Version int       `json:"version"`
	}
)
