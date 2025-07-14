package crypto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	t.Parallel()

	type test struct {
		name     string
		password string
	}

	tests := []test{
		{
			name:     "valid password",
			password: "T3s|n9p@ssw0rd!",
		},
		{
			name:     "empty",
			password: "",
		},
		{
			name:     "long string",
			password: strings.Repeat("x", 50),
		},
	}

	for _, ttest := range tests {
		tt := ttest
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash, err := HashPassword(tt.password)
			require.NoError(t, err)
			assert.NotEmpty(t, hash)

			valid := CheckPasswordHash(tt.password, hash)
			assert.True(t, valid)
		})
	}
}
