package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathsMatch(t *testing.T) {
	t.Parallel()

	t.Run("different len", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		st := []Step{
			{Name: "test1", Index: -1},
			{Name: "test2", Index: -1},
			{Name: "test3", Index: -1},
		}

		res := PathsMatch(st, []string{"test1"})

		rq.False(res)
	})

	t.Run("different len: bigger", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		st := []Step{
			{Name: "test1", Index: -1},
			{Name: "test2", Index: -1},
			{Name: "test3", Index: -1},
		}

		res := PathsMatch(st, []string{"test1", "test2", "test3", "test4"})

		rq.False(res)
	})

	t.Run("different steps", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		st := []Step{
			{Name: "test1", Index: -1},
			{Name: "test2", Index: -1},
			{Name: "test3", Index: -1},
		}

		res := PathsMatch(st, []string{"test1", "test2", "test4"})

		rq.False(res)
	})

	t.Run("same", func(t *testing.T) {
		t.Parallel()
		rq := require.New(t)

		st := []Step{
			{Name: "test1", Index: -1},
			{Name: "test2", Index: -1},
			{Name: "test3", Index: -1},
		}

		res := PathsMatch(st, []string{"test1", "test2", "test3"})

		rq.True(res)
	})
}
