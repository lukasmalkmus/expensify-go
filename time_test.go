package expensify

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTime_MarshalJSON_ZeroValue(t *testing.T) {
	tm := Time{}

	b, err := json.Marshal(tm)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	assert.EqualValues(t, "null", b)
}

func TestTime_MarshalJSON(t *testing.T) {
	now := NewTime(time.Now())

	b, err := json.Marshal(now)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	str := fmt.Sprintf("\"%d-%02d-%02d\"", now.Year(), now.Month(), now.Day())
	assert.EqualValues(t, str, string(b))
}

func TestTime_UnmarshalJSON_ZeroValue(t *testing.T) {
	var tm Time
	err := json.Unmarshal([]byte("null"), &tm)
	require.NoError(t, err)

	assert.True(t, tm.IsZero())
}

func TestTime_UnmarshalJSON(t *testing.T) {
	now := time.Now()
	str := fmt.Sprintf("\"%d-%02d-%02d\"", now.Year(), now.Month(), now.Day())

	var tm Time
	err := json.Unmarshal([]byte(str), &tm)
	require.NoError(t, err)

	assert.EqualValues(t, strings.Trim(str, "\""), tm.Format(layoutISO))
}
