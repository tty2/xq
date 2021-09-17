package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainsString(t *testing.T) {
	t.Parallel()

	t.Run("true", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		rq.True(ContainsString([]string{"1", "2", "3"}, "2"))
	})

	t.Run("false", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		rq.False(ContainsString([]string{"1", "2", "3"}, "4"))
	})
}
