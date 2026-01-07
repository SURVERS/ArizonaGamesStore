package handlers

import (
	"regexp"
	"strings"
)

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	if strings.Contains(errStr, "23505") {
		return true
	}

	return false
}

func IsValidNickname(nickname string) (bool, string) {
	if nickname == "" {
		return false, "Никнейм не может быть пустым"
	}

	if len(nickname) < 3 || len(nickname) > 20 {
		return false, "Никнейм должен быть от 3 до 20 символов"
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, nickname)
	if !matched {
		return false, "Никнейм может содержать только латинские буквы, цифры и символы _ . -"
	}

	return true, ""
}
