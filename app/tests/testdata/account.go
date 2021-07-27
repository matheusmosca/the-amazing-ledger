package testdata

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateAccount() string {
	return "liability.clients.available." + strings.ReplaceAll(uuid.New().String(), "-", "_") + ".*"
}
