package main

import (
	"bufio"

	"github.com/tty2/xq/internal/formatter"
	"github.com/tty2/xq/internal/processor"
)

const (
	indentItemSize = 2
)

type prc interface {
	Process(r *bufio.Reader) chan string
}

func getProcessor(q query) (prc, error) {
	if len(q.path) == 0 {
		return formatter.New(indentItemSize)
	}

	return processor.New(q.path, q.attribute, q.searchType)
}
