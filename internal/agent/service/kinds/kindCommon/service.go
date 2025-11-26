package kindCommon

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	crypro "github.com/dontagr/pandora/pkg/crypto"
)

const oneMBInBytes = 1024 * 1024

type Kind struct {
	NeedData bool
	NeedFile bool
	Hasher   *crypro.CManager
}

func (m *Kind) EncryptData(data string) (string, error) {
	encrypt, err := m.Hasher.Encrypt(bytes.NewBufferString(data))
	if err != nil {
		return "", fmt.Errorf("EncryptData: %v", err)
	}

	return encrypt.String(), nil
}

func (m *Kind) DecryptData(data string) (string, error) {
	decrypt, err := m.Hasher.Decrypt([]byte(data))
	if err != nil {
		return "", fmt.Errorf("DecryptData: %v", err)
	}

	return string(decrypt), nil
}

func (m *Kind) IsDataRequired() bool {
	return m.NeedData
}

func (m *Kind) IsFileRequired() bool {
	return m.NeedFile
}

func (m *Kind) CommonValidateData(data string) error {
	if data == "" {
		return fmt.Errorf("json data required")
	}

	if !isStringLessThanOneMB(data) {
		return fmt.Errorf("json size more than 1 MB sory but it's too much")
	}

	return nil
}

func (m *Kind) ValidateLabel(label string) error {
	if label == "" {
		return fmt.Errorf("label is required")
	}

	if len(label) > 255 {
		return fmt.Errorf("label len more than 256")
	}

	return nil
}

func (m *Kind) ValidateMeta(meta string) error {
	if meta != "" && !isStringLessThanOneMB(meta) {
		return fmt.Errorf("meta len more than one MB")
	}

	return nil
}

func (m *Kind) ValidateFile(file string) error {
	if file == "" {
		return fmt.Errorf("file path required")
	}

	err := isFileReadable(file)
	if err != nil {
		return err
	}

	err = isFileSizeLessThanOneMB(file)
	if err != nil {
		return err
	}

	return nil
}
func isFileReadable(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filename)
		}

		if os.IsPermission(err) {
			return fmt.Errorf("permission denied for file: %s", filename)
		}

		return fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()

	return nil
}

func isFileSizeLessThanOneMB(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	sizeInBytes := info.Size()
	if sizeInBytes < oneMBInBytes {
		return nil
	}

	return fmt.Errorf("file size: %.2f MB, it's more than 1 MB", bytesToMegabytes(sizeInBytes))
}

func bytesToMegabytes(bytes int64) float64 {
	return float64(bytes) / oneMBInBytes
}

func isStringLessThanOneMB(s string) bool {
	return len(s) < oneMBInBytes
}
func (m *Kind) StructToMap(data interface{}) map[string]string {
	result := make(map[string]string)
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	if val.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < val.NumField(); i++ {
		key := strings.ToLower(typ.Field(i).Name)
		value := val.Field(i).Interface()

		result[key] = fmt.Sprintf("%v", value)
	}

	return result
}
