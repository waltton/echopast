package whois

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWhois(t *testing.T) {
	testCases := []struct {
		name     string
		registry string
	}{
		{
			name:     "1.1.1.1-iana",
			registry: RegistryIANA,
		},
		{
			name:     "13.13.13.13-arin",
			registry: RegistryARIN,
		},
		{
			name:     "1.1.1.1-apnic",
			registry: RegistryAPNIC,
		},
		{
			name:     "186.186.186.186-lacnic",
			registry: RegistryLACNIC,
		},
		{
			name:     "102.0.0.0-afrinic",
			registry: RegistryAFRINIC,
		},
		{
			name:     "104.167.16.0-ripe",
			registry: RegistryRIPE,
		},
	}

	for _, tc := range testCases {
		raw, err := os.ReadFile("./data/" + tc.name + ".txt")
		require.NoError(t, err)

		f, err := os.Open("./data/" + tc.name + ".json")
		require.NoError(t, err)

		ew := Whois{}
		err = json.NewDecoder(f).Decode(&ew)
		require.NoError(t, err)

		result, err := parseWhois(string(raw))
		require.NoError(t, err)

		if len(ew) != len(result) {
			fmt.Printf("result: %+v\n", result)
		}

		assert.Equal(t, len(ew), len(result))
		for objc := range ew {
			for k := range ew[objc] {
				val, ok := result[objc][k]
				require.True(t, ok, "no value found for objc: %d, key: %s, object: %+v", objc, k, result[objc])
				assert.Equal(t, ew[objc][k], val)
			}
		}

		if tc.registry != "" {
			assert.Equal(t, tc.registry, result.Registry())
		}
	}
}
