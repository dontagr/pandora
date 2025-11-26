package models

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	UrlUserSignUP  = "api/user/register"
	UrlUserLogin   = "api/user/login"
	UrlTest        = "api/store/test"
	UrlStoreSave   = "api/store/save"
	UrlStoreLoad   = "api/store/load"
	UrlStoreDelete = "api/store/delete"

	KindAuth   = "auth"
	KindText   = "text"
	KindBinary = "binary"
	KindCard   = "card"
)

type (
	ResponceStoreLoad struct {
		Kind       string            `json:"kind"`
		Label      string            `json:"label"`
		Version    int               `json:"version"`
		Data       string            `json:"data"`
		Meta       string            `json:"meta"`
		DT         time.Time         `json:"dt"`
		ReveryData map[string]string `json:"-"`
	}
	RequestStoreLoad struct {
		Kind  string `json:"kind"`
		Label string `json:"label"`
	}
	RequestStoreSave struct {
		Kind  string `json:"kind"`
		Data  any    `json:"data"`
		Meta  string `json:"meta"`
		Label string `json:"label"`
	}
	RequestUser struct {
		Login    string `json:"login" validate:"required,alphanum|email"`
		Password string `json:"password" validate:"required"`
	}
	CommonResponce struct {
		Body   []byte
		Meta   CommonMeta
		Status int
	}
	CommonMeta struct {
		Authorization string
	}
	RequestKindAuth struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	RequestKindCard struct {
		Number string `json:"number"`
		Date   string `json:"date"`
		Cvv    string `json:"cvv"`
	}
	RequestKindText struct {
		Text string `json:"text"`
	}

	RequestKindBinary struct {
		Text string `json:"text"`
	}
)

func (r RequestStoreSave) GetData() (string, error) {
	marshal, err := json.Marshal(r.Data)
	if err != nil {
		return "", fmt.Errorf("error marshaling data: %w", err)
	}

	return string(marshal), nil
}
