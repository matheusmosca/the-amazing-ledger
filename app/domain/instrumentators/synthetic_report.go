package instrumentators

import (
	"context"
	"fmt"
	"time"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (lp *LedgerInstrumentator) GettingSyntheticReport(ctx context.Context, account vos.Account, startTime time.Time, endTime time.Time) {
	lp.Log(ctx, fmt.Sprintf("getting synthetic report: %s/%s-%s", account.Value(), startTime.String(), endTime.String()))
}

func (lp *LedgerInstrumentator) GotSyntheticReport(ctx context.Context, report vos.SyntheticReport) {
	lp.Log(ctx, fmt.Sprintf("got synthetic report: num_accounts = %d, total_credit = %d, total_debit = %d", len(report.Results), report.TotalCredit, report.TotalDebit))
}
