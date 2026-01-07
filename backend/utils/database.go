package utils

import "strings"

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	if strings.Contains(errStr, "23505") {
		return true
	}

	if strings.Contains(strings.ToLower(errStr), "duplicate key") {
		return true
	}

	if strings.Contains(strings.ToLower(errStr), "unique constraint") {
		return true
	}

	return false
}
