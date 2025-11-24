package interfaces

type (
	Kind interface {
		GetName() string
		IsDataRequired() bool
		IsFileRequired() bool
		ValidateData(data string) error
		ValidateLabel(data string) error
		ValidateFile(data string) error
		GetStruct(data string, file string) (any, error)
		Encrypt(data any) (any, error)
	}
)
