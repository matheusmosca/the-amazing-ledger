package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

const syntheticReportQuery = `
select 
	subpath(account, 0, $1),
	coalesce(SUM(CASE operation WHEN %d THEN amount ELSE 0::bigint END),0::bigint) AS creditSum, 
	coalesce(SUM(CASE operation WHEN %d THEN amount ELSE 0::bigint END),0::bigint) AS debitSum 
from 
	entry 
where 
	account ~ $2
and 
	created_at >= $3 and created_at < $4 
group by 1;
`

func (r *LedgerRepository) GetSyntheticReport(ctx context.Context, query vos.Account, level int, startTime time.Time, endTime time.Time) (*vos.SyntheticReport, error) {
	const operation = "Repository.GetSyntheticReport"

	sqlQuery, params := buildQueryAndParams(query, level, startTime, endTime)

	defer r.pb.MonitorDataSegment(ctx, collection, operation, sqlQuery).End()
	rows, errQuery := r.db.Query(
		ctx,
		sqlQuery,
		params...,
	)

	if errQuery != nil {
		if errors.Is(errQuery, pgx.ErrNoRows) {
			return &vos.SyntheticReport{}, nil
		}
		return nil, errQuery
	}

	defer rows.Close()

	results := []vos.AccountResult{}
	var totalCredit int64
	var totalDebit int64

	for rows.Next() {
		var accStr string
		var credit int64
		var debit int64

		err := rows.Scan(
			&accStr,
			&credit,
			&debit,
		)

		if err != nil {
			return nil, err
		}

		account, err := vos.NewAnalyticalAccount(accStr)
		if err != nil {
			return nil, err
		}

		path := vos.AccountResult{
			Account: account,
			Credit:  credit,
			Debit:   debit,
		}

		results = append(results, path)

		totalCredit = totalCredit + credit
		totalDebit = totalDebit + debit
	}

	errNext := rows.Err()
	if errNext != nil {
		return nil, errNext
	}

	if results == nil || len(results) < 1 {
		return &vos.SyntheticReport{}, nil
	}

	syntheticReport, errEntity := vos.NewSyntheticReport(totalCredit, totalDebit, results)
	if errEntity != nil {
		return nil, errEntity
	}

	return syntheticReport, nil
}

func buildQueryAndParams(query vos.Account, level int, startTime time.Time, endTime time.Time) (string, []interface{}) {
	sqlQuery := syntheticReportQuery
	sqlQuery = fmt.Sprintf(sqlQuery, vos.CreditOperation, vos.DebitOperation)

	params := make([]interface{}, 0)
	params = append(params, strconv.Itoa(level), query.Value(), startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	return sqlQuery, params
}
