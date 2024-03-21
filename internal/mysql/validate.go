package mysql

import (
	"fmt"
	"regexp"
)

func sanitizeDatabaseName(databaseName string) (string, error) {
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	if !validNameRegex.MatchString(databaseName) {
		return "", fmt.Errorf("invalid database name")
	}

	if len(databaseName) > 64 {
		return "", fmt.Errorf("database name too long")
	}

	return databaseName, nil
}

func sanitizeUsername(databaseName string) (string, error) {
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	if !validNameRegex.MatchString(databaseName) {
		return "", fmt.Errorf("invalid username")
	}

	if len(databaseName) > 32 {
		return "", fmt.Errorf("username name too long")
	}

	return databaseName, nil
}

func sanitizePassword(password string) (string, error) {
	if len(password) < 8 {
		return "", fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 32 {
		return "", fmt.Errorf("password too long")
	}

	return password, nil
}
