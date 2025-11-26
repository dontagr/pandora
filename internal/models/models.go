package models

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	UrlUserSignUP  = "api/user/register"
	UrlUserLogin   = "api/user/login"
	UrlStoreSave   = "api/store/save"
	UrlStoreList   = "api/store/list"
	UrlStoreLoad   = "api/store/load"
	UrlStoreDelete = "api/store/delete"
	UrlStoreSync   = "api/store/sync"

	KindAuth   = "auth"
	KindText   = "text"
	KindBinary = "binary"
	KindCard   = "card"
)

type (
	ResponceStoreList struct {
		List []StoreLite `json:"list"`
	}
	StoreLite struct {
		DT      time.Time `json:"dt"`
		Label   string    `json:"label"`
		Version int       `json:"version"`
	}
	ResponceSync struct {
		Version int `json:"version"`
	}
	ResponceStoreLoad struct {
		DT         time.Time         `json:"dt"`
		ReveryData map[string]string `json:"-"`
		Kind       string            `json:"kind"`
		Label      string            `json:"label"`
		Data       string            `json:"data"`
		Meta       string            `json:"meta"`
		Version    int               `json:"version"`
	}
	RequestSyncList struct {
		List []SyncNode `json:"list"`
	}
	SyncNode struct {
		Kind    string `json:"kind"`
		Label   string `json:"label"`
		Data    string `json:"data"`
		Meta    string `json:"meta"`
		Version int    `json:"version"`
	}
	RequestStoreLoad struct {
		Kind  string `json:"kind"`
		Label string `json:"label"`
	}
	RequestStoreList struct {
		Kind string `json:"kind"`
	}
	RequestStoreSave struct {
		Kind    string `json:"kind"`
		Data    any    `json:"data"`
		Meta    string `json:"meta"`
		Label   string `json:"label"`
		Version int    `json:"version"`
	}
	RequestUser struct {
		Login    string `json:"login" validate:"required,alphanum|email"`
		Password string `json:"password" validate:"required"`
	}
	CommonResponce struct {
		Meta   CommonMeta
		Body   []byte
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
