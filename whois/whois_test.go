package whois

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {
	result, err := Lookup("131.131.131.131")
	require.NoError(t, err)

	// w, err := parseWhoisIANA(result)
	// require.NoError(t, err)

	// fmt.Printf("w: %+v", w)

	fmt.Println("result", result)
}

func TestIsZeroString(t *testing.T) {
	testCases := []struct {
		value   string
		expects bool
	}{
		{
			value:   "",
			expects: true,
		},
		{
			value:   "asd",
			expects: false,
		},
		{
			value:   string([]byte{0, 0, 0, 0}),
			expects: true,
		},
	}

	for _, tc := range testCases {
		result := isZeroString(tc.value)
		assert.Equal(t, tc.expects, result)
	}
}
