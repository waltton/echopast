package whois

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWhois(t *testing.T) {
	testCases := []struct {
		name     string
		registry string
		refer    string
		country  string
	}{
		{
			name:     "1.1.1.1-iana",
			registry: RegistryIANA,
			refer:    "whois.apnic.net",
		},
		{
			name:     "13.13.13.13-arin",
			registry: RegistryARIN,
			country:  "US",
		},
		{
			name:     "1.1.1.1-apnic",
			registry: RegistryAPNIC,
			country:  "AU",
		},
		{
			name:     "186.186.186.186-lacnic",
			registry: RegistryLACNIC,
			country:  "VE",
		},
		{
			name:     "102.0.0.0-afrinic",
			registry: RegistryAFRINIC,
			country:  "KE",
		},
		{
			name:     "104.167.16.0-ripe",
			registry: RegistryRIPE,
			country:  "DE",
		},
		{
			name:     "131.131.131.131-iana",
			registry: RegistryIANA,
			refer:    "whois.arin.net",
		},
		{
			name:     "131.131.131.131-arin",
			registry: RegistryARIN,
			country:  "US",
		},
		{
			name:     "163.123.142.153-apnic",
			registry: RegistryAPNIC,
			country:  "AU",
		},
		{
			name:     "3.81.245.94-arin",
			registry: RegistryARIN,
			country:  "US",
		},
		{
			name:     "152.89.196.211-arin",
			registry: RegistryARIN,
			country:  "NL",
		},
		{
			name:     "188.166.125.106-ripe",
			registry: RegistryRIPE,
			country:  "NL",
		},
		// {
		// 	name:     "51.79.29.48-ripe",
		// 	registry: RegistryARIN,
		// 	country:  "CA",
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := os.ReadFile("./data/" + tc.name + ".txt")
			require.NoError(t, err)

			f, err := os.Open("./data/" + tc.name + ".json")
			require.NoError(t, err)

			ew := Whois{}
			err = json.NewDecoder(f).Decode(&ew.Data)
			require.NoError(t, err)

			result, err := parseWhois(string(raw))
			require.NoError(t, err)

			assert.Equal(t, len(ew.Data), len(result.Data))
			for objc := range ew.Data {
				for k := range ew.Data[objc] {
					val, ok := result.Data[objc][k]
					require.True(t, ok, "no value found for objc: %d, key: %s, object: %+v", objc, k, result.Data[objc])
					assert.Equal(t, ew.Data[objc][k], val)
				}
			}

			if tc.registry != "" {
				assert.Equal(t, tc.registry, result.Registry())
			}

			if tc.refer != "" {
				assert.Equal(t, tc.refer, result.Refer())
			}

			if tc.country != "" {
				assert.Equal(t, tc.country, result.Country())
			}
		})
	}
}
