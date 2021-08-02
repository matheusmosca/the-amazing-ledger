package usecases

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func (l *LedgerUseCase) GetAccountBalance(ctx context.Context, account vos.Account) (vos.AccountBalance, error) {
	var (
		accountBalance vos.AccountBalance
		err            error
	)

	switch account.Type() {
	case vos.Analytic:
		accountBalance, err = l.repository.GetAnalyticAccountBalance(ctx, account)
	case vos.Synthetic:
		accountBalance, err = l.repository.GetSyntheticAccountBalance(ctx, account)
	default:
		err = app.ErrInvalidAccountType
	}

	if err != nil {
		return vos.AccountBalance{}, fmt.Errorf("failed to get account balance: %w", err)
	}

	return accountBalance, nil
}
