package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnalyticalAccountBalance(t *testing.T) {
	account, err := NewAnalyticalAccount("liability.clients.available.user_1.block")
	assert.NoError(t, err)

	accountBalance := NewAnalyticalAccountBalance(account, Version(3), 100, 50)

	assert.Equal(t, AccountBalance{
		Account:        account,
		CurrentVersion: Version(3),
		TotalCredit:    100,
		TotalDebit:     50,
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
		TotalCredit:    0,
		TotalDebit:     0,
		Balance:        50,
	}, accountBalance)
}
