package testdata

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateAccount() string {
	return "liability.clients.available." + strings.ReplaceAll(uuid.New().String(), "-", "_") + ".*"
}

func GenerateInvalidAccount() string {
	return "liability.clients." + strings.ReplaceAll(uuid.New().String(), "-", "_") + ".*"
}
