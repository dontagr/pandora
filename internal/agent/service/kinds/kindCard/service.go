package kindAuth

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dontagr/pandora/internal/agent/service/kinds/factory"
	"github.com/dontagr/pandora/internal/agent/service/kinds/kindCommon"
	"github.com/dontagr/pandora/internal/models"
	crypro "github.com/dontagr/pandora/pkg/crypto"
)

type Kind struct {
	kindCommon.Kind
	name string
}

func RegisterKind(mf *factory.KindFactory, hasher *crypro.CManager) {
	mf.SetKind(&Kind{
		name: models.KindCard,
		Kind: kindCommon.Kind{
			NeedData: true,
			Hasher:   hasher,
		},
	})
}

func (m *Kind) GetName() string {
	return m.name
}

func (m *Kind) ValidateData(data string) error {
	err := m.CommonValidateData(data)
	if err != nil {
		return err
	}

	return nil
}

func (m *Kind) GetStruct(data string, _ string) (any, error) {
	var req models.RequestKindCard

	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		return nil, fmt.Errorf("error with deserializing JSON: %v", err)
	}

	if req.Number == "" && req.Date == "" && req.Cvv == "" {
		return nil, fmt.Errorf("json doesn't contain Number or Date or Cvv")
	}

	if req.Number != "" && isValidCardNumber(req.Number) {
		return nil, fmt.Errorf("card number invalid")
	}

	if req.Date != "" && !isValidDate(req.Date) {
		return nil, fmt.Errorf("card date invalid")
	}

	if req.Cvv != "" && (len(req.Cvv) != 3 || isDigitOnly(req.Cvv)) {
		return nil, fmt.Errorf("card cvv invalid")
	}

	return &req, nil
}

func isDigitOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isValidCardNumber(number string) bool {
	number = strings.ReplaceAll(number, " ", "")

	if len(number) >= 13 && len(number) <= 19 && isDigitOnly(number) {
		return true
	}

	return false
}

func isValidDate(dateStr string) bool {
	re := regexp.MustCompile(`^(\d{1,2})/(\d{2})$`)
	matches := re.FindStringSubmatch(dateStr)
	if matches == nil {
		return false
	}

	month, err := strconv.Atoi(matches[1])
	if err != nil || month < 1 || month > 12 {
		return false
	}
	year, err := strconv.Atoi(matches[2])
	if err != nil {
		return false
	}
	year = year + 2000

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	if year > currentYear || (year == currentYear && time.Month(month) >= currentMonth) {
		return true
	}

	return false
}

func (m *Kind) Encrypt(data any) (any, error) {
	var err error
	req := data.(*models.RequestKindCard)

	if req.Number != "" {
		req.Number, err = m.EncryptData(req.Number)
		if err != nil {
			return nil, err
		}
	}

	if req.Date != "" {
		req.Date, err = m.EncryptData(req.Date)
		if err != nil {
			return nil, err
		}
	}

	if req.Cvv != "" {
		req.Cvv, err = m.EncryptData(req.Cvv)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}
