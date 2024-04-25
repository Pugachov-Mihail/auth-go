package auth_validate

import (
	"fmt"
	"log/slog"
	"regexp"
)

type AuthPersonalDataValidate struct {
	log *slog.Logger
	Validate
}

type Validate interface {
	validateEmail(email string) (bool, error)
	ValidateLenValuesString(value string, maxLen int) (bool, error)
}

func validateEmail(email string) (bool, error) {
	result, err := regexp.MatchString(`^[A-z0-9]*@[A-z0-9-]*\.[A-z]{2,4}$`, email)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки почты")
	}
	if result {
		return true, nil
	}
	return false, nil
}

func ValidateLenValuesString(value string, maxLen int) bool {
	if len(value) <= maxLen {
		return true
	}
	return false
}
