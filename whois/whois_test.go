package whois

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

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
		line, err := rd.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		require.NoError(t, err)

		i++
		n := 2
		if i < n || i > n {
			continue
		}

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
