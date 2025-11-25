package models

const (
	UrlUserSignUP = "api/user/register"
	UrlUserLogin  = "api/user/login"
	UrlTest       = "api/store/test"
	UrlStoreSave  = "api/store/save"

	KindAuth   = "auth"
	KindText   = "text"
	KindBinary = "binary"
	KindCard   = "card"
)

type (
	RequestStoreSave struct {
		Kind  string `json:"kind"`
		Data  any    `json:"data"`
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
