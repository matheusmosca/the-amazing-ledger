package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestLedgerRepository_GetSyntheticReportSuccess(t *testing.T) {
	r := NewLedgerRepository(pgDocker.DB, logrus.New())
	ctx := context.Background()

	accountBase := "liability.assets"
	accountBaseEmpty := "liability.income"

	account, err := vos.NewAccount(accountBase + ".account11")
	assert.NoError(t, err)

	e1 := createEntry(t, vos.DebitOperation, account.Value(), vos.NextAccountVersion, 100)
	e2 := createEntry(t, vos.CreditOperation, account.Value(), vos.IgnoreAccountVersion, 100)

	tx := createTransaction(t, ctx, r, e1, e2)

	e3 := createEntry(t, vos.DebitOperation, account.Value(), vos.IgnoreAccountVersion, 100)
	e4 := createEntry(t, vos.CreditOperation, account.Value(), vos.IgnoreAccountVersion, 100)

	tx2 := createTransaction(t, ctx, r, e3, e4)

	e5 := createEntry(t, vos.DebitOperation, account.Value(), vos.IgnoreAccountVersion, 100)
	e6 := createEntry(t, vos.CreditOperation, account.Value(), vos.IgnoreAccountVersion, 100)

	tx3 := createTransaction(t, ctx, r, e5, e6)

	e7 := createEntry(t, vos.DebitOperation, account.Value(), vos.IgnoreAccountVersion, 100)
	e8 := createEntry(t, vos.CreditOperation, account.Value(), vos.IgnoreAccountVersion, 100)

	tx4 := createTransaction(t, ctx, r, e7, e8)

	e9 := createEntry(t, vos.DebitOperation, account.Value(), vos.IgnoreAccountVersion, 100)
	e10 := createEntry(t, vos.CreditOperation, account.Value(), vos.IgnoreAccountVersion, 100)

	tx5 := createTransaction(t, ctx, r, e9, e10)

	testCases := []struct {
		name        string
		query       string
		level       int
		startDate   time.Time
		endDate     time.Time
		transaction entities.Transaction
		report      vos.SyntheticReport
	}{
		{
			name:        "01 - should get a result because the account path is correct",
			query:       account.Value(),
			level:       3,
			startDate:   time.Now().UTC(),
			endDate:     time.Now().UTC().Add(time.Hour * 1),
			transaction: tx,
			report: vos.SyntheticReport{
				TotalCredit: 500,
				TotalDebit:  500,
			},
		},
		{
			name:        "02 - should get a result because the path is correct w/ query",
			query:       accountBase + ".*",
			level:       3,
			startDate:   time.Now().UTC(),
			endDate:     time.Now().UTC().Add(time.Hour * 1),
			transaction: tx2,
			report: vos.SyntheticReport{
				TotalCredit: 500,
				TotalDebit:  500,
			},
		},
		{
			name:        "03 - should not get a result because there was no data inserted",
			query:       accountBaseEmpty + ".*",
			level:       3,
			startDate:   time.Now().UTC(),
			endDate:     time.Now().UTC().Add(time.Hour * 1),
			transaction: tx3,
			report:      vos.SyntheticReport{},
		},
		{
			name:        "04 - should not get a result because the account is wrong",
			query:       accountBase + ".omni",
			level:       3,
			startDate:   time.Now().UTC(),
			endDate:     time.Now().UTC().Add(time.Hour * 1),
			transaction: tx4,
			report:      vos.SyntheticReport{},
		},
		{
			name:        "05 - should not get a result because there are no results for the chosen period",
			query:       accountBase + ".*",
			level:       3,
			startDate:   time.Now().UTC(),
			endDate:     time.Now().UTC(),
			transaction: tx5,
			report:      vos.SyntheticReport{},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			query, err := vos.NewAccount(tt.query)
			assert.NoError(t, err)

			got, err := r.GetSyntheticReport(ctx, query, tt.level, tt.startDate, tt.endDate)
			assert.NoError(t, err)
			assert.Equal(t, tt.report.TotalCredit, got.TotalCredit)
			assert.Equal(t, tt.report.TotalDebit, got.TotalDebit)
		})
	}
}
