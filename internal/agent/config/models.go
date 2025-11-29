// Package config предоставляет структуры и функции для управления конфигурацией приложения.
// Он определяет структуру конфигурации, поддерживающую криптографию, логирование, работу с базой данных и сетевые настройки.
package config

var (
	// version хранит текущую версию приложения.
	version = "0.0.1"
	// bildDT хранит дату и время компиляции.
	bildDT = ""
)

// Config представляет собой основную структуру конфигурации приложения.
// Она включает в себя конфигурацию для криптографии, логирования, работы с базой данных и сетевых настроек.
type Config struct {
	Crypto     Crypto   `json:"Crypto"`
	Log        Logging  `json:"Logging"`
	DataBase   DataBase `json:"DataBase"`
	Version    string
	BildDT     string
	HTTPServer HTTPServer `json:"HttpServing"`
	Transport  Transport  `json:"Transport"`
}

// Crypto содержит параметры конфигурации для работы с криптографией.
type Crypto struct {
	PublicKey  string `json:"PublicKey" validate:"required"`
	PrivateKey string `json:"PrivateKey" validate:"required"`
}

// Logging содержит настройки логирования.
type Logging struct {
	LogLevel string `json:"LogLevel" validate:"required"`
}

// Transport содержит параметры конфигурации для управления сетевыми параметрами.
type Transport struct {
	WaitForRetry int `json:"WaitForRetry" validate:"required"`
}

// HTTPServer содержит параметры конфигурации для работы HTTP-сервера.
type HTTPServer struct {
	Host string `json:"Host" validate:"required"`
	Gzip bool   `json:"Gzip"`
}

// DataBase содержит параметры конфигурации для работы с базой данных.
type DataBase struct {
	Path string `json:"Path" validate:"required"`
}

// NewConfig создает и возвращает новый экземпляр структуры Config.
// Этот метод инициализирует значения версии и даты компиляции.
func NewConfig() *Config {
	cnf := Config{}
	cnf.Version = version
	cnf.BildDT = bildDT

	return &cnf
}
