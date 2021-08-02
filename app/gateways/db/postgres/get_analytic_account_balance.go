package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/instrumentation/newrelic"
)

const getAccountBalanceQuery = `
select
	total_balance,
    version
from
    get_analytic_account_balance($1)
;
`

func (r LedgerRepository) GetAnalyticAccountBalance(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
	const operation = "Repository.GetAnalyticAccountBalance"

	defer newrelic.NewDatastoreSegment(ctx, collection, operation, getAccountBalanceQuery).End()

	var balance int
	var version int64

	err := r.db.QueryRow(ctx, getAccountBalanceQuery, account.Value()).Scan(
		&balance,
		&version,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return vos.AccountBalance{}, fmt.Errorf("failed to get account balance: %w", err)
		}

		if pgErr.Code == pgerrcode.NoDataFound {
			return vos.AccountBalance{}, app.ErrAccountNotFound
		}

		return vos.AccountBalance{}, fmt.Errorf("failed to get account balance: %w", pgErr)
	}

	return vos.NewAnalyticAccountBalance(
		account,
		vos.Version(version),
		balance,
	), nil
}
