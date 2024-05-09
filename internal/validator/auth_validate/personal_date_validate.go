package auth_validate

import (
	"fmt"
	"regexp"
)

func ValidateEmail(email string) (bool, error) {
	result, err := regexp.MatchString(`^[\w.-]+@[A-z0-9-]*\.[A-z]{2,4}$`, email)
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
