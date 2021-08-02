package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnalyticAccountBalance(t *testing.T) {
	account, err := NewAnalyticAccount("liability.clients.available.user_1.block")
	assert.NoError(t, err)

	accountBalance := NewAnalyticAccountBalance(account, Version(3), 50)

	assert.Equal(t, AccountBalance{
		Account:        account,
		CurrentVersion: Version(3),
		Balance:        50,
	}, accountBalance)
}

func TestNewSyntheticAccountBalance(t *testing.T) {
	account, err := NewAccount("liability.clients.available.user_1.*")
	assert.NoError(t, err)

	accountBalance := NewSyntheticAccountBalance(account, 50)

	assert.Equal(t, AccountBalance{
		Account:        account,
		CurrentVersion: IgnoreAccountVersion,
		Balance:        50,
	}, accountBalance)
}
