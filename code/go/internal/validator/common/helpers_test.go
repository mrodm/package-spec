package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsDefinedWarningsAsErrors(t *testing.T) {
	cases := []struct {
		name        string
		envVarValue string
		expected    bool
	}{
		{"true", "true", true},
		{"false", "false", false},
		{"other", "other", false},
		{"empty", "", false},
		{"undefined", nil, false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if test.envVarValue != nil {
				if err := os.Setenv(EnvVarWarningsAsErrors, test.envVarValue); err != nil {
					require.NoError(t, err)
				}
			}
			value := IsDefinedWarningsAsErrors()
			assert.Equal(t, test.expected, value)
		})
	}
	if err := os.Unsetenv(EnvVarWarningsAsErrors); err != nil {
		require.NoError(t, err)
	}
}
