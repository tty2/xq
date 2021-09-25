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

func TestTagSetName(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("<tagname>"),
		}
		rq.Equal("", tg.Name)

		err := tg.SetName()
		rq.NoError(err)
		rq.Equal("tagname", tg.Name)
	})

	t.Run("close tag", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("</tagname>"),
		}
		rq.Equal("", tg.Name)

		err := tg.SetName()
		rq.NoError(err)
		rq.Equal("tagname", tg.Name)
	})

	t.Run("empty tag", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("<tagname />"),
		}
		rq.Equal("", tg.Name)

		err := tg.SetName()
		rq.NoError(err)
		rq.Equal("tagname", tg.Name)
	})

	t.Run("with attribute", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("<tagname attr='value'>"),
		}
		rq.Equal("", tg.Name)

		err := tg.SetName()
		rq.NoError(err)
		rq.Equal("tagname", tg.Name)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("tagname attr='value'>"),
		}
		rq.Equal("", tg.Name)

		err := tg.SetName()
		rq.Error(err)
	})
}

func TestTagSetNameAndAttributes(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("<tagname attr='value'>"),
		}
		rq.Equal("", tg.Name)
		rq.Len(tg.Attributes, 0)

		err := tg.SetNameAndAttributes()
		rq.NoError(err)
		rq.Equal("tagname", tg.Name)
		rq.Len(tg.Attributes, 1)
		v, ok := tg.Attributes["attr"]
		rq.True(ok)
		rq.Equal("value", v)
	})

	t.Run("ok: several attributes", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte(`<tagname attr="value" data="datavalue" number="42">`),
		}
		rq.Equal("", tg.Name)
		rq.Len(tg.Attributes, 0)

		err := tg.SetNameAndAttributes()
		rq.NoError(err)
		rq.Equal("tagname", tg.Name)
		rq.Len(tg.Attributes, 3)
		v, ok := tg.Attributes["attr"]
		rq.True(ok)
		rq.Equal("value", v)
		d, ok := tg.Attributes["data"]
		rq.True(ok)
		rq.Equal("datavalue", d)
		n, ok := tg.Attributes["number"]
		rq.True(ok)
		rq.Equal("42", n)
		rq.Equal("datavalue", d)
		f, ok := tg.Attributes["notfound"]
		rq.False(ok)
		rq.Equal("", f)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("tagname attr='value'>"),
		}
		rq.Equal("", tg.Name)

		err := tg.SetNameAndAttributes()
		rq.Error(err)
	})

	t.Run("close tag", func(t *testing.T) {
		t.Parallel()

		tg := Tag{
			Bytes: []byte("</tagname>"),
		}
		rq.Equal("", tg.Name)

		err := tg.SetNameAndAttributes()
		rq.NoError(err)
		rq.Equal("tagname", tg.Name)
		rq.Len(tg.Attributes, 0)
	})
}
