package symbol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsQuote(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		rq := require.New(t)

		rq.True(IsQuote('\''))
		rq.True(IsQuote('"'))
		rq.False(IsQuote('\n'))
		rq.False(IsQuote('\t'))
		rq.False(IsQuote('a'))
		rq.False(IsQuote('A'))
	})
}

func Test_IsSpace(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		rq := require.New(t)

		rq.True(IsSpace(' '))
		rq.True(IsSpace('\n'))
		rq.True(IsSpace('\t'))
		rq.False(IsSpace('1'))
		rq.False(IsSpace('a'))
		rq.False(IsSpace('A'))
		rq.False(IsSpace('z'))
	})
}
