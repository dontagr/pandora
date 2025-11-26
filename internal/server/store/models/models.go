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
	Secret struct {
		Kind    string    `json:"kind"`
		Label   string    `json:"label"`
		Version int       `json:"version"`
		Data    string    `json:"data"`
		Meta    string    `json:"meta"`
		DT      time.Time `json:"dt"`
	}
)
