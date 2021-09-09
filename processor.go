package main

import (
	"bufio"

	"github.com/tty2/xq/internal/processors/attributes"
	"github.com/tty2/xq/internal/processors/data"
	"github.com/tty2/xq/internal/processors/tags"
)

const (
	indentItemSize = 2
)

type processor interface {
	Process(r *bufio.Reader) error
}

func getProcessor(q query) (processor, error) {
	if len(q.path) == 0 {
		return data.NewProcessor(indentItemSize)
	} else if q.attribute != "" {
		return attributes.NewProcessor(
			q.path,
			q.attribute,
		)
	}

	return tags.NewProcessor(q.path)
}
