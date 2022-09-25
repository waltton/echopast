package whois

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func TestLookupFromFile(t *testing.T) {
	f, err := os.Open("./data/ips.txt")
	require.NoError(t, err)

	var i int
	rd := bufio.NewReader(f)
	for {
		i++

		if i > 1 {
			return
		}

		line, err := rd.ReadString('\n')
		require.NoError(t, err)

		line = strings.TrimSpace(line)

		t.Run(line, func(t *testing.T) {
			result, err := Lookup(line)
			require.NoError(t, err)

			w, err := parseWhois(result)
			require.NoError(t, err)

			fmt.Println(line)
			fmt.Println("Country", w.Country())
			fmt.Println("Refer", w.Refer())
			fmt.Println("Registry", w.Registry())
			fmt.Println("-------")

		})

	}
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
