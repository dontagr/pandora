// Пакет models предоставляет основные модели и структуры данных,
// используемые в приложении для управления пользователями и хранилищем.
package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Константы, определяющие URL-адреса API для различных операций пользователя и хранилища.
const (
	UrlUserSignUP  = "api/user/register" // URL для регистрации пользователя.
	UrlUserLogin   = "api/user/login"    // URL для входа пользователя.
	UrlStoreSave   = "api/store/save"    // URL для сохранения данных в хранилище.
	UrlStoreList   = "api/store/list"    // URL для получения списка данных из хранилища.
	UrlStoreLoad   = "api/store/load"    // URL для загрузки данных из хранилища.
	UrlStoreDelete = "api/store/delete"  // URL для удаления данных из хранилища.
	UrlStoreSync   = "api/store/sync"    // URL для синхронизации данных.
)

// Константы, определяющие виды данных в хранилище.
const (
	KindAuth   = "auth"   // Указывает на данные типа аутентификации.
	KindText   = "text"   // Указывает на текстовые данные.
	KindBinary = "binary" // Указывает на бинарные данные.
	KindCard   = "card"   // Указывает на данные банковской карты.
)

type (
	// ResponceStoreList представляет ответ на запрос списка сущностей в хранилище.
	ResponceStoreList struct {
		List []StoreLite `json:"list"`
	}
	// StoreLite представляет "легкое" представление сущности в хранилище,
	// включающее в себя метку, дату и версию данных.
	StoreLite struct {
		DT      time.Time `json:"dt"`
		Label   string    `json:"label"`
		Version int       `json:"version"`
	}
	// ResponceSync представляет ответ на запрос синхронизации.
	ResponceSync struct {
		Version int `json:"version"`
	}
	// ResponceSyncList представляет ответ на запрос списка сущностей для синхронизации.
	ResponceSyncList struct {
		List []*ResponceStoreLoad `json:"list"`
	}
	// ResponceSyncListSave представляет данные, которые должны быть сохранены в процессе синхронизации.
	ResponceSyncListSave struct {
		List []*RequestStoreSave `json:"list"`
	}
	// ResponceStoreLoad представляет ответ на запрос загрузки сущности из хранилища.
	ResponceStoreLoad struct {
		DT           time.Time         `json:"dt"`
		RecoveryData map[string]string `json:"-"`
		Kind         string            `json:"kind"`
		Label        string            `json:"label"`
		Data         string            `json:"data"`
		Meta         string            `json:"meta"`
		Version      int               `json:"version"`
	}
	// RequestSyncList представляет запрос на синхронизацию данных.
	RequestSyncList struct {
		List []*RequestStoreSave `json:"list"`
	}
	// SyncNode представляет ноду синхронизации.
	SyncNode struct {
		Kind    string `json:"kind"`
		Label   string `json:"label"`
		Data    string `json:"data"`
		Meta    string `json:"meta"`
		Version int    `json:"version"`
	}
	// RequestStoreLoad представляет собой структуру запроса для загрузки хранилища.
	RequestStoreLoad struct {
		Kind  string `json:"kind"`
		Label string `json:"label"`
	}
	// RequestStoreList представляет запрос на получение списка сущностей.
	RequestStoreList struct {
		Kind string `json:"kind"`
	}
	// RequestStoreSave представляет собой структуру запроса для сохранения хранилища.
	RequestStoreSave struct {
		Kind    string `json:"kind"`
		Data    any    `json:"data"`
		Meta    string `json:"meta"`
		Label   string `json:"label"`
		Version int    `json:"version"`
	}
	// RequestUser представляет запрос на действия с пользователем,
	// включающие в себя логин и пароль.
	RequestUser struct {
		Login    string `json:"login" validate:"required,alphanum|email"`
		Password string `json:"password" validate:"required"`
	}
	// CommonResponce представляет общий ответ, содержаший основную информацию.
	CommonResponce struct {
		Meta   CommonMeta
		Body   []byte
		Status int
	}
	// CommonMeta содержит метаинформацию о запросе, например, токен авторизации.
	CommonMeta struct {
		Authorization string
	}
	// RequestKindAuth представляет запрос для аутентификации.
	RequestKindAuth struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	// RequestKindCard представляет данные о банковской карте.
	RequestKindCard struct {
		Number string `json:"number"`
		Date   string `json:"date"`
		Cvv    string `json:"cvv"`
	}
	// RequestKindText представляет текстовые данные.
	RequestKindText struct {
		Text string `json:"text"`
	}
	// RequestKindBinary представляет бинарные данные.
	RequestKindBinary struct {
		Text string `json:"text"`
	}
)

// GetData возвращает данные структуры RequestStoreSave в виде строки JSON.
func (r RequestStoreSave) GetData() (string, error) {
	marshal, err := json.Marshal(r.Data)
	if err != nil {
		return "", fmt.Errorf("error marshaling data: %w", err)
	}

	return string(marshal), nil
}
