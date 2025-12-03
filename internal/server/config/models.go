package config

type Config struct {
	Log        Logging    `json:"Logging"`
	HTTPServer HTTPServer `json:"HttpServing"`
	DataBase   DataBase   `json:"DataBase"`
	Security   Security   `json:"Security"`
}

type Logging struct {
	LogLevel string `json:"LogLevel" validate:"required"`
}

type DataBase struct {
	DatabaseDsn string `json:"DatabaseDsn" validate:"required"`
}

type HTTPServer struct {
	BindAddress string `json:"BindAddress" validate:"required"`
}

type Security struct {
	Key string `json:"key" validate:"required"`
}
