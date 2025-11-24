package user

import (
	"fmt"
	"unicode"

	"github.com/spf13/cobra"

	"github.com/dontagr/pandora/internal/models"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (us *Service) GetLoginPassword(cmd *cobra.Command) (*models.RequestUser, error) {
	login, _ := cmd.Flags().GetString("login")
	password, _ := cmd.Flags().GetString("password")
	if login == "" {
		return nil, fmt.Errorf("login is required")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if !isValidPassword(password) {
		return nil, fmt.Errorf("password is invalid")
	}

	return &models.RequestUser{Login: login, Password: password}, nil
}

func isValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	var hasDigit, hasLower, hasUpper bool
	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		}
	}

	return hasDigit && hasLower && hasUpper
}
