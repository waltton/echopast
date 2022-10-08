package whois

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {

	result, err := Lookup("131.131.131.131")
	require.NoError(t, err)

	fmt.Println("result", result)
}

func TestLookupFromFile(t *testing.T) {
	workers := 4
	var wg sync.WaitGroup

	ips := make(chan string, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for ip := range ips {
				t.Run(ip, func(t *testing.T) {
					w, err := Lookup(ip)
					require.NoError(t, err)

					assert.NotEmpty(t, w.Country())

				})
			}
		}()
	}

	f, err := os.Open("./data/ips.txt")
	require.NoError(t, err)

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		require.NoError(t, err)

		ips <- strings.TrimSpace(line)
	}

	close(ips)
	wg.Wait()
}
