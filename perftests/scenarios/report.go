package scenarios

import "fmt"

// Scenario 02
//
// Type: banking/output
// Description: get synthetic accout reports after an insert benchmark
// Steps:
// 	- execute Scenario 01
// 	- query for a synthetic report

type ReportScenario struct {
	method       string
	accountQuery string
}

func NewReportScenario(account string) ReportScenario {
	s := ReportScenario{}
	s.method = "ledger.LedgerService.GetSyntheticReport"
	s.accountQuery = createAccountQuery(account)
	return s
}

func createAccountQuery(account string) string {
	return fmt.Sprintf(`{"account": "%s", "start_date": "%s", "end_date": "%s"}`, account+".*", "2020-01-01T15:04:05.999999999Z", "2030-01-02T15:04:05.999999999Z")
}

func (s ReportScenario) GetMethod() string {
	return s.method
}

func (s ReportScenario) GetJSON() string {
	return s.accountQuery
}
