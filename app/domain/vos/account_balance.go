package vos

type AccountBalance struct {
	Account        Account
	CurrentVersion Version
	Balance        int
}

func NewAnalyticAccountBalance(account Account, version Version, balance int) AccountBalance {
	return AccountBalance{
		Account:        account,
		CurrentVersion: version,
		Balance:        balance,
	}
}

func NewSyntheticAccountBalance(account Account, balance int) AccountBalance {
	return AccountBalance{
		Account:        account,
		CurrentVersion: IgnoreAccountVersion,
		Balance:        balance,
	}
}
