package scenarios

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Scenario 01
//
// Type: banking/input
// Description: create accounts and transactions
// Steps:
// 	- we create a "company" fake account (assets.aaa.bbb.sa)
// 	- we create 1k fake client accounts (liability.aaa.bbb.UUID)
// 	- for each client account above, we create a transaction that transfers money from the client account to the company account.
// 	- these transactions are executed N times.

type LedgerScenario struct {
	method         string
	companyAccount string
	clientAccounts []string
	transactions   string
}

func NewLedgerScenario(totalClientAccounts, totalRequests int) LedgerScenario {
	s := LedgerScenario{}
	s.method = "ledger.LedgerService.CreateTransaction"
	s.CreateAccounts(totalClientAccounts, totalRequests)
	s.CreateTransactions()
	return s
}

func (s *LedgerScenario) CreateAccounts(totalClientAccounts, totalRequests int) {
	// Create the Company account.
	s.companyAccount = "asset.aaa.bbb.sa"

	// Create N client accounts.
	s.clientAccounts = make([]string, totalClientAccounts)
	for i := 0; i < totalClientAccounts; i++ {
		uuid := uuid.New().String()
		uuid = strings.ReplaceAll(uuid, "-", "_")
		s.clientAccounts[i] = "liability.aaa.bbb." + uuid
	}
}

func (s *LedgerScenario) CreateTransactions() {
	var series strings.Builder
	series.WriteString("[")

	for _, ca := range s.clientAccounts {
		// 1 transaction that transfer money from clientaccoutn to company account
		series.WriteString(transferTransaction(s.companyAccount, ca) + ",")
	}

	s.transactions = series.String()
	s.transactions = s.transactions[:len(s.transactions)-1] + "]"
}

func (s LedgerScenario) GetJSON() string {
	return s.transactions
}

func (s LedgerScenario) GetMethod() string {
	return s.method
}

func (s LedgerScenario) GetRandomClientAccount() string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(s.clientAccounts))))
	if err != nil {
		panic(err)
	}
	return s.clientAccounts[n.Int64()]
}

func transferTransaction(companyAccount string, targetAccount string) string {
	e1 := entryAsString("{{newUUID}}", targetAccount, 0, "OPERATION_DEBIT", 20000)
	e2 := entryAsString("{{newUUID}}", companyAccount, 0, "OPERATION_CREDIT", 20000)
	tr := fmt.Sprintf(`{"id":"{{newUUID}}", "competence_date":"%v", "company":"%s","event":"1", "entries":[%s,%s]}`, time.Now().Format(time.RFC3339), targetAccount, e1, e2)
	return tr
}

func entryAsString(id string, accountID string, expectedVersion int, operation string, amount int) string {
	return fmt.Sprintf(`{"id":"%s","account":"%s", "expected_version": %d, "operation": "%s", "amount": %d}`, id, accountID, expectedVersion, operation, amount)
}
