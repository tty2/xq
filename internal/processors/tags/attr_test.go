package tags

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPickAttributesNames(t *testing.T) {
	t.Parallel()

	t.Run("nil: empty tag", func(t *testing.T) {
		t.Parallel()

		tag := []byte{}
		res := pickAttributesNames(tag)

		rq := require.New(t)
		rq.Nil(res)
	})

	t.Run("nil: first byte is not an open bracket", func(t *testing.T) {
		t.Parallel()

		tag := []byte("tagname")
		res := pickAttributesNames(tag)

		rq := require.New(t)
		rq.Nil(res)
	})

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		tag := []byte("<tagname attr1='value1' attr2='value2' attr3='value3'>")
		res := pickAttributesNames(tag)

		rq := require.New(t)
		rq.Len(res, 3)
		rq.Equal("attr1", res[0])
		rq.Equal("attr2", res[1])
		rq.Equal("attr3", res[2])
	})

	t.Run("ok: other quotes", func(t *testing.T) {
		t.Parallel()

		tag := []byte(`<tagname attr1="value1" attr2="value2" attr3="value3">`)
		res := pickAttributesNames(tag)

		rq := require.New(t)
		rq.Len(res, 3)
		rq.Equal("attr1", res[0])
		rq.Equal("attr2", res[1])
		rq.Equal("attr3", res[2])
	})
}

func TestPickAttributeValue(t *testing.T) {
	t.Parallel()

	rq := require.New(t)

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		v, err := pickAttributeValue("attr", []byte("<tagname attr='value' attr2='value2'>"))
		rq.NoError(err)
		rq.Equal("value", v)
	})

	t.Run("ok: double quotes", func(t *testing.T) {
		t.Parallel()

		v, err := pickAttributeValue("attr", []byte(`<tagname attr="value" attr2="value2">`))
		rq.NoError(err)
		rq.Equal("value", v)
	})

	t.Run("err: no open bracket", func(t *testing.T) {
		t.Parallel()

		_, err := pickAttributeValue("attr", []byte(`tagname attr="value" attr2="value2">`))
		rq.Error(err)
	})

	t.Run("err: no close bracket", func(t *testing.T) {
		t.Parallel()

		_, err := pickAttributeValue("attr", []byte(`<tagname attr="value" attr2="value2"`))
		rq.Error(err)
	})

	t.Run("err: too short", func(t *testing.T) {
		t.Parallel()

		_, err := pickAttributeValue("attr", []byte(`<>"`))
		rq.Error(err)
	})

	t.Run("err: didn't found attribute", func(t *testing.T) {
		t.Parallel()

		_, err := pickAttributeValue("attr", []byte(`<tagname attr1="value" attr2="value2">`))
		rq.Error(err)
	})
}
