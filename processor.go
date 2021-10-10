package main

import (
	"bufio"

	"github.com/tty2/xq/internal/processors/data"
	"github.com/tty2/xq/internal/processors/tags"
)

const (
	indentItemSize = 2
)

type processor interface {
	Process(r *bufio.Reader) chan string
}

func getProcessor(q query) (processor, error) {
	if len(q.path) == 0 {
		return data.NewProcessor(indentItemSize)
	}

	return tags.NewProcessor(q.path, q.attribute, q.searchType)
}
