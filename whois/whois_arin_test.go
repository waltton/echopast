package whois

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWhoisARIN(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{name: "13.13.13.13-arin"},
	}

	for _, tc := range testCases {
		raw, err := os.ReadFile("./data/" + tc.name + ".txt")
		require.NoError(t, err)

		f, err := os.Open("./data/" + tc.name + ".json")
		require.NoError(t, err)

		ew := WhoisARIN{}
		err = json.NewDecoder(f).Decode(&ew)
		require.NoError(t, err)

		result, err := parseWhoisARIN(string(raw))
		require.NoError(t, err)

		assert.Equal(t, len(ew), len(result))
	}
}
