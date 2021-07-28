package vos

type AccountBalance struct {
	Account        Account
	CurrentVersion Version
	TotalCredit    int
	TotalDebit     int
	Balance        int
}

func NewAnalyticalAccountBalance(account Account, version Version, totalCredit, totalDebit int) AccountBalance {
	return AccountBalance{
		Account:        account,
		CurrentVersion: version,
		TotalCredit:    totalCredit,
		TotalDebit:     totalDebit,
		Balance:        totalCredit - totalDebit,
	}
}

func NewSyntheticAccountBalance(account Account, balance int) AccountBalance {
	return AccountBalance{
		Account:        account,
		CurrentVersion: IgnoreAccountVersion,
		TotalCredit:    0,
		TotalDebit:     0,
		Balance:        balance,
	}
}
