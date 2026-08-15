package exec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRegistryCredentials(t *testing.T) {
	tests := []struct {
		name     string
		raw      []string
		expected int
		server   string
	}{
		{
			name:     "parses standard registry entry",
			raw:      []string{"user:pass:registry.example.com"},
			expected: 1,
			server:   "registry.example.com",
		},
		{
			name:     "parses registry entry with port",
			raw:      []string{"user:pass:registry.example.com:5050"},
			expected: 1,
			server:   "registry.example.com:5050",
		},
		{
			name:     "parses registry entry with multiple colons in server",
			raw:      []string{"user:pass:[2001:db8::1]:5000"},
			expected: 1,
			server:   "[2001:db8::1]:5000",
		},
		{
			name:     "skips malformed entry",
			raw:      []string{"malformed"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registries := ParseRegistryCredentials(tt.raw)
			require.Len(t, registries, tt.expected)
			if tt.expected > 0 {
				require.Equal(t, "user", registries[0].Username)
				require.Equal(t, "pass", registries[0].Password)
				require.Equal(t, tt.server, registries[0].ServerAddress)
			}
		})
	}
}
