package main

import (
	"fmt"
	"os"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"github.com/stone-co/the-amazing-ledger/perftests/scenarios"
)

const (
	totalClientAccounts = 1000
	totalRequests       = 15000
	concurrency         = 20
)

// we are going to execute multiple scenarios
// each of them will have it's own benchmark results
func main() {
	for _, s := range getScenarios() {

		fmt.Printf("> running scenario: %s\n", s.GetMethod())
		fmt.Printf(">> data preview: %.80s\n\n", s.GetJSON())

		report, err := runner.Run(
			s.GetMethod(),
			"0.0.0.0:3000",
			runner.WithProtoset("./protoimage.bin"),
			runner.WithDataFromJSON(s.GetJSON()),
			runner.WithInsecure(true),
			runner.WithTotalRequests(uint(totalRequests)),
			runner.WithConcurrency(uint(concurrency)),
		)

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		printer := printer.ReportPrinter{
			Out:    os.Stdout,
			Report: report,
		}

		printer.Print("summary")
	}
}

// getting scenarios
// each one has it own data and you can create a scenario using
// data from another one (same database)
func getScenarios() []scenarios.Scenario {
	ls := scenarios.NewLedgerScenario(totalClientAccounts, totalRequests)
	rs := scenarios.NewReportScenario(ls.GetRandomClientAccount())

	return []scenarios.Scenario{
		ls,
		rs,
	}
}
