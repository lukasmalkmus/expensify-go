package expensify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpenseService_Create(t *testing.T) {
	exp := &Expense{
		Merchant: "Apple Inc.",
		Created:  mustTimeParse(t, layoutISO, "2020-09-01"),
		Amount:   99,
		Currency: "EUR",
		ReportID: 1234,
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		require.NoError(t, r.ParseForm())
		payload := r.PostFormValue("requestJobDescription")

		var req map[string]interface{}
		err := json.Unmarshal([]byte(payload), &req)
		require.NoError(t, err)

		assert.Len(t, req, 3)
		assert.Contains(t, req, "type")
		assert.Contains(t, req, "credentials")
		if assert.Contains(t, req, "inputSettings") {
			if assert.Contains(t, req["inputSettings"], "transactionList") {
				tl := req["inputSettings"].(map[string]interface{})["transactionList"].([]interface{})
				tl1 := tl[0].(map[string]interface{})
				assert.EqualValues(t, "Apple Inc.", tl1["merchant"])
				assert.EqualValues(t, "2020-09-01", tl1["created"])
				assert.EqualValues(t, 99, tl1["amount"])
				assert.EqualValues(t, "EUR", tl1["currency"])
				assert.EqualValues(t, 1234, tl1["reportID"])
			}
		}

		_, err = fmt.Fprint(w, `{
			"responseCode" : 200,
			"transactionList" : [
				{
					"merchant" : "Apple Inc.",
					"created" : "2020-09-01",
					"amount" : 99,
					"transactionID" : "82827382377292",
					"currency" : "EUR",
					"reportID" : 1234
				}
			]
		}`)
		require.NoError(t, err)
	}

	client, teardown := setup(t, hf)
	defer teardown()

	res, err := client.Expense.Create(context.Background(), "dev@example.com", []*Expense{exp})
	require.NoError(t, err)

	if assert.NotEmpty(t, res) {
		assert.EqualValues(t, exp.Merchant, res[0].Merchant)
		assert.EqualValues(t, exp.Created, res[0].Created)
		assert.EqualValues(t, exp.Amount, res[0].Amount)
		assert.EqualValues(t, "82827382377292", res[0].TransactionID)
		assert.EqualValues(t, exp.Currency, res[0].Currency)
		assert.EqualValues(t, exp.ReportID, res[0].ReportID)
	}
}

func mustTimeParse(t *testing.T, layout, value string) Time {
	ts, err := time.Parse(layout, value)
	if err != nil {
		require.NoError(t, err)
	}
	return NewTime(ts)
}
