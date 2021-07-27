package vos

import "github.com/stone-co/the-amazing-ledger/app"

// TODO: improve struct name(Common Language)
type AccountResult struct {
	Account Account
	Credit  int64
	Debit   int64
}

type SyntheticReport struct {
	TotalCredit int64
	TotalDebit  int64
	Results     []AccountResult
}

func NewSyntheticReport(totalCredit, totalDebit int64, accounts []AccountResult) (*SyntheticReport, error) {
	if accounts == nil || len(accounts) < 1 {
		return nil, app.ErrInvalidSyntheticReportStructure
	}

	return &SyntheticReport{
		TotalCredit: totalCredit,
		TotalDebit:  totalDebit,
		Results:     accounts,
	}, nil
}
