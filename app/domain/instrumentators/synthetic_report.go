package instrumentators

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (lp *LedgerInstrumentator) GettingSyntheticReport(ctx context.Context, account vos.Account, startTime time.Time, endTime time.Time) {
	zerolog.Ctx(ctx).Info().
		Str("account", account.Value()).
		Str("start_time", startTime.String()).
		Str("end_time", endTime.String()).
		Msg("getting synthetic report")
}

func (lp *LedgerInstrumentator) GotSyntheticReport(ctx context.Context, report vos.SyntheticReport) {
	zerolog.Ctx(ctx).Info().
		Int("accounts_total", len(report.Results)).
		Msg("got synthetic report")
}
