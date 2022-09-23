package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildParams(t *testing.T) {
	testCases := []struct {
		cols, rows int
		expects    string
	}{
		{
			cols:    1,
			rows:    1,
			expects: "($1)",
		},
		{
			cols:    2,
			rows:    1,
			expects: "($1,$2)",
		},
		{
			cols:    1,
			rows:    2,
			expects: "($1),($2)",
		},
		{
			cols:    2,
			rows:    2,
			expects: "($1,$2),($3,$4)",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expects, buildParams(tc.cols, tc.rows), "cols: %d, rows: %d", tc.cols, tc.rows)
	}
}
