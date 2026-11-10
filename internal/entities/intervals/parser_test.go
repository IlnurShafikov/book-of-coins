package intervals

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIntervals(t *testing.T) {
	errIs := func(expErr error) func(t *testing.T, gotErr error) {
		return func(t *testing.T, gotErr error) {
			t.Helper()
			require.ErrorIs(t, gotErr, expErr)
		}
	}

	tests := []struct {
		name     string
		data     [][]string
		exp      Intervals
		checkErr func(t *testing.T, gotErr error)
	}{
		{
			name:     "not enough value",
			data:     [][]string{{"1"}},
			exp:      nil,
			checkErr: errIs(ErrInvalidData),
		},
		{
			name: "invalid key",
			data: [][]string{{"abc", "100", "200"}},
			exp:  nil,
			checkErr: func(t *testing.T, gotErr error) {
				t.Helper()
				require.Error(t, gotErr)
				require.Contains(t, gotErr.Error(), "invalid ID")

				var expErr *strconv.NumError
				require.ErrorAs(t, gotErr, &expErr, "invalid ID")
			},
		},
		{
			name: "invalid min value",
			data: [][]string{{"1", "xyz", "130"}},
			exp:  nil,
			checkErr: func(t *testing.T, gotErr error) {
				t.Helper()
				require.Error(t, gotErr)
				require.Contains(t, gotErr.Error(), "invalid Min")

				var expErr *strconv.NumError
				require.ErrorAs(t, gotErr, &expErr, "invalid Min")
			},
		},
		{
			name: "invalid max value",
			data: [][]string{{"1", "130", "xyz"}},
			exp:  nil,
			checkErr: func(t *testing.T, gotErr error) {
				t.Helper()
				require.Error(t, gotErr)
				require.Contains(t, gotErr.Error(), "invalid Max")

				var expErr *strconv.NumError
				require.ErrorAs(t, gotErr, &expErr, "invalid Max")
			},
		},
		{
			name: "parse data successful",
			data: [][]string{
				{"1", "0", "99"},
				{"2", "100", "200"},
			},
			exp: Intervals{
				Interval{ID: 1, Min: 0, Max: 99},
				Interval{ID: 2, Min: 100, Max: 200},
			},
			checkErr: errIs(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Parse(tt.data)
			tt.checkErr(t, err)
			require.Equal(t, tt.exp, resp)
		})
	}
}
