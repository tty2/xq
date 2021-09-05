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

type searchQuery struct {
	query query
}

type parser struct {
	CurrentTag     tag
	Data           []byte
	IndentItemSize int
	Indentation    int
	InsideTag      bool // semaphore that shows if we read data inside a tag
	SkipData       bool
	searchQuery    searchQuery
}

type tag struct {
	Name     string
	String   string
	Bytes    []byte
	Brackets int
}

type attribute struct {
	Name        []byte
	Value       []byte
	Quote       byte
	NextIsQuote bool
	InsideValue bool
}

func newParser(q query) parser {
	return parser{
		IndentItemSize: indentItemSize,
		Data:           []byte{},
		CurrentTag: tag{
			Bytes: []byte{},
		},
		searchQuery: searchQuery{
			query: q,
		},
	}
}

func (p *parser) getProcessor() processor {
	if len(p.searchQuery.query.path) == 0 {
		return data.NewProcessor(indentItemSize)
	} else if p.searchQuery.query.attribute != "" {
		return attributes.NewProcessor(
			p.searchQuery.query.path,
			p.searchQuery.query.attribute,
		)
	}

	return tags.NewProcessor(p.searchQuery.query.path)
}
