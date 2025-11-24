package config

var (
	version = "0.0.1"
	bildDT  = ""
)

type Config struct {
	Log        Logging    `json:"Logging"`
	Transport  Transport  `json:"Transport"`
	HTTPServer HTTPServer `json:"HttpServing"`
	DataBase   DataBase   `json:"DataBase"`
	Crypto     Crypto     `json:"Crypto"`
	Version    string
	BildDT     string
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
