package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
	"github.com/stone-co/the-amazing-ledger/app/tests/testdata"
	"github.com/stone-co/the-amazing-ledger/app/tests/testenv"
	"github.com/stone-co/the-amazing-ledger/app/tests/testseed"
	"github.com/stone-co/the-amazing-ledger/app/tests/testutils"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
)

type createTransactionBody struct {
	ID             string  `json:"id"`
	Entries        []entry `json:"entries"`
	CompetenceDate string  `json:"competence_date"`
	Company        string  `json:"company"`
	Event          uint32  `json:"event"`
}

type entry struct {
	ID              string `json:"id"`
	Account         string `json:"account"`
	ExpectedVersion int64  `json:"expected_version"`
	Operation       int8   `json:"operation"`
	Amount          int64  `json:"amount"`
}

func TestE2E_Gateway_CreateTransactionSuccess(t *testing.T) {
	t.Run("should create a transaction successfully", func(t *testing.T) {
		requestBody := createTransactionBody{
			ID: uuid.New().String(),
			Entries: []entry{
				{
					ID:              uuid.New().String(),
					Account:         testdata.GenerateAccountPath(),
					ExpectedVersion: vos.NextAccountVersion.AsInt64(),
					Operation:       int8(vos.DebitOperation),
					Amount:          100,
				},
				{
					ID:              uuid.New().String(),
					Account:         testdata.GenerateAccountPath(),
					ExpectedVersion: vos.NextAccountVersion.AsInt64(),
					Operation:       int8(vos.CreditOperation),
					Amount:          100,
				},
			},
			Company:        "abc",
			CompetenceDate: time.Now().Format(time.RFC3339),
			Event:          1,
		}

		req := testutils.NewRequest(t, http.MethodPost, testenv.GatewayServer+"/api/v1/transactions", requestBody)
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		resp, err := client.Do(req)
		assert.NoError(t, err)

		defer resp.Body.Close()

		var body bytes.Buffer

		_, err = io.Copy(&body, resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, `{}`, body.String())

		tests.TruncateTables(context.Background(), testenv.DB, "entry", "account_version")
	})
}

func TestE2E_Gateway_CreateTransactionFailure(t *testing.T) {
	type (
		responseBody struct {
			Code    int
			Message string
		}
		wants struct {
			status int
			body   responseBody
		}
	)

	testCases := []struct {
		name     string
		seedRepo func(t *testing.T) entities.Transaction
		request  createTransactionBody
		wants    wants
	}{
		{
			name: "should return an error if id is invalid",
			request: createTransactionBody{
				ID: "invalid UUID",
				Entries: []entry{
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.DebitOperation),
						Amount:          100,
					},
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.CreditOperation),
						Amount:          100,
					},
				},
				Company:        "abc",
				CompetenceDate: time.Now().Format(time.RFC3339),
				Event:          1,
			},
			wants: wants{
				status: http.StatusBadRequest,
				body: responseBody{
					Code:    3,
					Message: "error parsing transaction id",
				},
			},
		},
		{
			name: "should return an error if entry id is invalid",
			request: createTransactionBody{
				ID: uuid.New().String(),
				Entries: []entry{
					{
						ID:              "invalid entry id",
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.DebitOperation),
						Amount:          100,
					},
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.CreditOperation),
						Amount:          100,
					},
				},
				Company:        "abc",
				CompetenceDate: time.Now().Format(time.RFC3339),
				Event:          1,
			},
			wants: wants{
				status: http.StatusBadRequest,
				body: responseBody{
					Code:    3,
					Message: "error parsing entry id",
				},
			},
		},
		{
			name: "should return an error if operation is invalid",
			request: createTransactionBody{
				ID: uuid.New().String(),
				Entries: []entry{
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(proto.Operation_OPERATION_UNSPECIFIED),
						Amount:          100,
					},
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.CreditOperation),
						Amount:          100,
					},
				},
				Company:        "abc",
				CompetenceDate: time.Now().Format(time.RFC3339),
				Event:          1,
			},
			wants: wants{
				status: http.StatusBadRequest,
				body: responseBody{
					Code:    3,
					Message: "invalid operation",
				},
			},
		},
		{
			name: "should return an error if amount is invalid",
			request: createTransactionBody{
				ID: uuid.New().String(),
				Entries: []entry{
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.DebitOperation),
						Amount:          -100,
					},
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.CreditOperation),
						Amount:          100,
					},
				},
				Company:        "abc",
				CompetenceDate: time.Now().Format(time.RFC3339),
				Event:          1,
			},
			wants: wants{
				status: http.StatusBadRequest,
				body: responseBody{
					Code:    3,
					Message: "invalid amount",
				},
			},
		},
		{
			name: "should return an error if number of entries is less than two",
			request: createTransactionBody{
				ID: uuid.New().String(),
				Entries: []entry{
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.DebitOperation),
						Amount:          100,
					},
				},
				Company:        "abc",
				CompetenceDate: time.Now().Format(time.RFC3339),
				Event:          1,
			},
			wants: wants{
				status: http.StatusConflict,
				body: responseBody{
					Code:    10,
					Message: "invalid entries number",
				},
			},
		},
		{
			name: "should return if competence date is in the future",
			request: createTransactionBody{
				ID: uuid.New().String(),
				Entries: []entry{
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.DebitOperation),
						Amount:          100,
					},
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.CreditOperation),
						Amount:          100,
					},
				},
				Company:        "abc",
				CompetenceDate: time.Now().Add(1 * time.Minute).Format(time.RFC3339),
				Event:          1,
			},
			wants: wants{
				status: http.StatusBadRequest,
				body: responseBody{
					Code:    3,
					Message: "competence date set to the future",
				},
			},
		},
		{
			name: "should return an error when occurs idempotency key violation",
			seedRepo: func(t *testing.T) entities.Transaction {
				e1 := testutils.CreateEntry(t, vos.DebitOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)
				e2 := testutils.CreateEntry(t, vos.CreditOperation, testdata.GenerateAccountPath(), vos.NextAccountVersion, 100)

				tx := testseed.CreateTransaction(t, e1, e2)

				return tx
			},
			request: createTransactionBody{
				ID: uuid.New().String(),
				Entries: []entry{
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.DebitOperation),
						Amount:          100,
					},
					{
						ID:              uuid.New().String(),
						Account:         testdata.GenerateAccountPath(),
						ExpectedVersion: vos.NextAccountVersion.AsInt64(),
						Operation:       int8(vos.CreditOperation),
						Amount:          100,
					},
				},
				Company:        "abc",
				CompetenceDate: time.Now().Format(time.RFC3339),
				Event:          1,
			},
			wants: wants{
				status: http.StatusBadRequest,
				body: responseBody{
					Code:    3,
					Message: "failed to create transaction: idempotency key violation",
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			request := tt.request

			if tt.seedRepo != nil {
				tx := tt.seedRepo(t)
				request.Entries[0].ID = tx.Entries[0].ID.String()

				defer tests.TruncateTables(context.Background(), testenv.DB, "entry", "account_version")
			}

			req := testutils.NewRequest(t, http.MethodPost, testenv.GatewayServer+"/api/v1/transactions", tt.request)
			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			resp, err := client.Do(req)
			assert.NoError(t, err)

			defer resp.Body.Close()

			var body responseBody

			err = json.NewDecoder(resp.Body).Decode(&body)
			assert.NoError(t, err)

			assert.Equal(t, tt.wants.status, resp.StatusCode)
			assert.Equal(t, tt.wants.body, body)
		})
	}
}
