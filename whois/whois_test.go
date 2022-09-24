package whois

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookup(t *testing.T) {
	result, err := Lookup("1.1.1.1")
	require.NoError(t, err)

	fmt.Println("result", result)
}
