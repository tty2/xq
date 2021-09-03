package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStep(t *testing.T) {
	t.Parallel()

	t.Run("clean tag", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)
		tg := "tag"

		res := getStep(tg)

		rq.Equal(tg, res.name)
		rq.Equal(-1, res.count)
	})

	t.Run("with index", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)
		tg := "tag_name[3]"

		res := getStep(tg)

		rq.Equal("tag_name", res.name)
		rq.Equal(3, res.count)
	})
}

func TestSeparateAttribute(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: "tag1.tag2#attr",
		}

		q.path = q.getPath()

		rq.Len(q.path, 2)

		res := q.separateAttribute()

		rq.Len(res, 2)
		rq.Equal("tag2", res[0])
		rq.Equal("attr", res[1])
	})
}

func TestGetPath(t *testing.T) {
	t.Parallel()

	t.Run("empty path", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: ".",
		}

		st := q.getPath()

		rq.Len(st, 0)
	})

	t.Run("without leading dot", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: "tag1.tag2.tag3",
		}

		st := q.getPath()

		rq.Len(st, 3)
		rq.Equal("tag1", st[0].name)
		rq.Equal(-1, st[0].count)
		rq.Equal("tag2", st[1].name)
		rq.Equal(-1, st[1].count)
		rq.Equal("tag3", st[2].name)
		rq.Equal(-1, st[2].count)
	})

	t.Run("leading dot", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: ".tag1.tag2.tag3",
		}

		st := q.getPath()

		rq.Len(st, 3)
		rq.Equal("tag1", st[0].name)
		rq.Equal(-1, st[0].count)
		rq.Equal("tag2", st[1].name)
		rq.Equal(-1, st[1].count)
		rq.Equal("tag3", st[2].name)
		rq.Equal(-1, st[2].count)
	})
}

func TestGetQuery(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		os.Args = []string{"xq"}

		q := getQuery()

		rq.Equal(".", q.request)
		rq.Equal(empty, q.target)
	})

	t.Run("path only", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		os.Args = []string{"xq", ".tag1.tag2"}

		q := getQuery()

		rq.Equal(".tag1.tag2", q.request)
		rq.Equal(_tags, q.target)
	})

	t.Run("path only", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		os.Args = []string{"xq", "tag", ".tag1.tag2"}

		q := getQuery()

		rq.Equal(".tag1.tag2", q.request)
		rq.Equal(_tags, q.target)
	})
}

func TestToTag(t *testing.T) {
	t.Parallel()

	t.Run("tag", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		rq.Equal(_tags, toTag("tag"))
	})

	t.Run("value", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		rq.Equal(tagValue, toTag("value"))
	})

	t.Run("attribute", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		rq.Equal(attr, toTag("attribute"))
	})

	t.Run("aValue", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		rq.Equal(attrValue, toTag("aValue"))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		rq.Equal(empty, toTag(""))
	})
}

func TestParseQuery(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			target:  _tags,
			request: ".",
		}

		q.parse()

		rq.Len(q.path, 0)
	})

	t.Run("_tags only", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			target:  _tags,
			request: ".tag1.tag2",
		}

		q.parse()

		rq.Len(q.path, 2)
		rq.Equal("tag1", q.path[0].name)
		rq.Equal("tag2", q.path[1].name)
	})

	t.Run("with attribute", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			target:  _tags,
			request: ".tag1.tag2#attr_name",
		}

		q.parse()

		rq.Len(q.path, 2)
		rq.Equal("tag1", q.path[0].name)
		rq.Equal("tag2", q.path[1].name)
		rq.Equal("attr_name", q.attribute)
	})
}
