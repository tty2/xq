package main

import (
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
