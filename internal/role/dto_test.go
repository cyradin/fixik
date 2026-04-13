package role

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		src      Permission
		expected []string
	}{
		{
			name:     "empty",
			src:      0,
			expected: nil,
		},
		{
			name:     "invalid permission",
			src:      1 << 63,
			expected: nil,
		},
		{
			name:     "incident get, status get",
			src:      IncidentGet | StatusGet,
			expected: []string{"INCIDENT_GET", "STATUS_GET"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.src.Codes()
			require.ElementsMatch(t, tt.expected, result)
		})
	}
}
