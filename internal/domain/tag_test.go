package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateTag(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("too short", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("<b"),
		}

		err := tg.Validate()
		rq.Error(err)
		rq.True(errors.Is(err, ErrTagShort))
	})

	t.Run("invalid start", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("tagName attr='value'>"),
		}

		err := tg.Validate()
		rq.Error(err)
		rq.True(errors.Is(err, ErrTagInvalidStart))
	})

	t.Run("invalid end", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("<tagName attr='value'"),
		}

		err := tg.Validate()
		rq.Error(err)
		rq.True(errors.Is(err, ErrTagInvalidEnd))
	})

	t.Run("invalid end", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("<tagName attr='value'>"),
		}

		err := tg.Validate()
		rq.NoError(err)
	})
}
