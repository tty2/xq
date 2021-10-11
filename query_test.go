package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tty2/xq/internal/domain"
)

func TestGetStep(t *testing.T) {
	t.Parallel()

	t.Run("clean tag", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)
		tg := "tag"

		res := getStep(tg)

		rq.Equal(tg, res.Name)
		rq.Equal(-1, res.Index)
	})

	t.Run("with index", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)
		tg := "tag_name[3]"

		res := getStep(tg)

		rq.Equal("tag_name", res.Name)
		rq.Equal(3, res.Index)
	})
}

func TestGetAttribute(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: "tag1.tag2#attr",
		}

		q.path = q.getPath()
		rq.Len(q.path, 2)

		rq.Len(q.path, 2)

		res := q.getAttribute()
		rq.Equal("attr", res)
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
		rq.Equal("tag1", st[0].Name)
		rq.Equal(-1, st[0].Index)
		rq.Equal("tag2", st[1].Name)
		rq.Equal(-1, st[1].Index)
		rq.Equal("tag3", st[2].Name)
		rq.Equal(-1, st[2].Index)
	})

	t.Run("leading dot", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: ".tag1.tag2.tag3",
		}

		st := q.getPath()

		rq.Len(st, 3)
		rq.Equal("tag1", st[0].Name)
		rq.Equal(-1, st[0].Index)
		rq.Equal("tag2", st[1].Name)
		rq.Equal(-1, st[1].Index)
		rq.Equal("tag3", st[2].Name)
		rq.Equal(-1, st[2].Index)
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
		rq.Empty(q.firstArg)
	})

	t.Run("path only", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		os.Args = []string{"xq", ".tag1.tag2"}

		q := getQuery()

		rq.Equal(".tag1.tag2", q.request)
		rq.Empty(q.firstArg)
	})

	t.Run("path only", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		os.Args = []string{"xq", "tags", ".tag1.tag2"}

		q := getQuery()

		rq.Equal(".tag1.tag2", q.request)
		rq.Equal("tags", q.firstArg)
	})

	t.Run("path only", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		os.Args = []string{"xq", "attr", ".tag1.tag2#val"}

		q := getQuery()

		rq.Equal(".tag1.tag2#val", q.request)
		rq.Equal("attr", q.firstArg)
	})
}

func TestParseQuery(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: ".",
		}

		q.parse()

		rq.Len(q.path, 0)
		rq.Equal(domain.TagValue, q.searchType)
	})

	t.Run("_tags only", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: ".tag1.tag2",
		}

		q.parse()

		rq.Len(q.path, 2)
		rq.Equal("tag1", q.path[0].Name)
		rq.Equal("tag2", q.path[1].Name)
		rq.Equal(domain.TagValue, q.searchType)
	})

	t.Run("_tags only: tag list", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request:  ".tag1.tag2",
			firstArg: "tags",
		}

		q.parse()

		rq.Len(q.path, 2)
		rq.Equal("tag1", q.path[0].Name)
		rq.Equal("tag2", q.path[1].Name)
		rq.Equal(domain.TagList, q.searchType)
	})

	t.Run("_tags only: attr list", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request:  ".tag1.tag2",
			firstArg: "attr",
		}

		q.parse()

		rq.Len(q.path, 2)
		rq.Equal("tag1", q.path[0].Name)
		rq.Equal("tag2", q.path[1].Name)
		rq.Equal(domain.AttrList, q.searchType)
	})

	t.Run("with attribute", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		q := query{
			request: ".tag1.tag2#attr_name",
		}

		q.parse()

		rq.Len(q.path, 2)
		rq.Equal("tag1", q.path[0].Name)
		rq.Equal("tag2", q.path[1].Name)
		rq.Equal("attr_name", q.attribute)
		rq.Equal(domain.AttrValue, q.searchType)
	})
}
