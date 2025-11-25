package models

import "time"

type (
	User struct {
		ID           int    `json:"id"`
		Login        string `json:"login"`
		PasswordHash string `json:"password"`
	}
	Secret struct {
		UserId  int       `json:"user_id"`
		Kind    string    `json:"kind"`
		Label   string    `json:"label"`
		Version int       `json:"version"`
		Data    string    `json:"data"`
		DT      time.Time `json:"dt"`
	}
)
