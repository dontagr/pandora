package config

var (
	version = "0.0.1"
	bildDT  = ""
)

type Config struct {
	Crypto     Crypto     `json:"Crypto"`
	Log        Logging    `json:"Logging"`
	HTTPServer HTTPServer `json:"HttpServing"`
	DataBase   DataBase   `json:"DataBase"`
	Version    string
	BildDT     string
	Transport  Transport `json:"Transport"`
}

type Crypto struct {
	PublicKey  string `json:"PublicKey" validate:"required"`
	PrivateKey string `json:"PrivateKey" validate:"required"`
}

type Logging struct {
	LogLevel string `json:"LogLevel" validate:"required"`
}

type Transport struct {
	WaitForRetry int `json:"WaitForRetry" validate:"required"`
}

type HTTPServer struct {
	Host string `json:"Host" validate:"required"`
}

type DataBase struct {
	Path string `json:"Path" validate:"required"`
}

func NewConfig() *Config {
	cnf := Config{}
	cnf.Version = version
	cnf.BildDT = bildDT

	return &cnf
}
