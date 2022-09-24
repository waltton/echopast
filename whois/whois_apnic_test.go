package whois

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWhoisAPNIC(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{name: "1.1.1.1-apnic"},
	}

	for _, tc := range testCases {
		raw, err := os.ReadFile("./data/" + tc.name + ".txt")
		require.NoError(t, err)

		f, err := os.Open("./data/" + tc.name + ".json")
		require.NoError(t, err)

		ew := WhoisAPNIC{}
		err = json.NewDecoder(f).Decode(&ew)
		require.NoError(t, err)

		result, err := parseWhoisAPNIC(string(raw))
		require.NoError(t, err)

		assert.Equal(t, len(ew), len(result))
		for objc := range ew {
			for k := range ew[objc] {
				val, ok := result[objc][k]
				require.True(t, ok, "no value found for objc: %d, key: %s, object: %+v", objc, k, result[objc])
				assert.Equal(t, ew[objc][k], val)
			}
		}
	}
}
