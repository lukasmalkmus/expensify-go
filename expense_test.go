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
			}
		}

		_, _ = fmt.Fprint(w, "{}")
	}

	client, teardown := setup(t, hf)
	defer teardown()

	_, err := client.Expense.Create(context.Background(), "dev@example.com", exp)
	require.NoError(t, err)
}

func TestUnmarshalTransactionList(t *testing.T) {
	var payload = `{
		"responseCode" : 200,
		"transactionList" : [
			{
				"amount" : 1234,
				"merchant" : "Name Of Merchant 1",
				"created" : "2016-01-01",
				"transactionID" : "6720309558248016",
				"currency" : "USD",
				"reportID":65343384
			},
			{
				"amount" : 2211,
				"merchant" : "Name Of Merchant 2",
				"created" : "2016-01-31",
				"transactionID" : "6720309558248017",
				"currency" : "CAD",
				"reportID":65343384
			}
		]
	}`

	res := new(CreateResponse)
	err := json.Unmarshal([]byte(payload), res)
	require.NoError(t, err)
	assert.Len(t, res.TransactionList, 2)
	assert.EqualValues(t, "6720309558248017", res.TransactionList[1].TransactionID)

	date, err := time.Parse(layoutISO, "2016-01-01")
	require.NoError(t, err)
	assert.EqualValues(t, NewTime(date), res.TransactionList[0].Created)
}

func mustTimeParse(t *testing.T, layout, value string) Time {
	ts, err := time.Parse(layout, value)
	if err != nil {
		require.NoError(t, err)
	}
	return NewTime(ts)
}
