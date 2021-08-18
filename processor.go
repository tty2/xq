package main

import (
	"bufio"
	"io"
)

const (
	closeBracket = '>'
	openBracket  = '<'

	minTagSize = 3 // minimum tag size can be 3. as example <b>

	red   = "\033[01;31m"
	green = "\033[01;32m"
	white = "\033[00m"

	indentItemSize = 2

	newLine        = 10 // '\n'
	carriageReturn = 13 // '\r'

	quote       = 39 // '
	doubleQuote = 34 // "

	space = 32
)

type searchQuery struct {
	// count int
	// print bool
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

func (p *parser) process(r *bufio.Reader) error {
	if len(p.searchQuery.query.path) == 0 {
		return p.fullProccess(r)
	}

	return nil
}

func (p *parser) fullProccess(r *bufio.Reader) error {
	buf := make([]byte, 0, 4*1024)

	for {
		n, err := r.Read(buf[:cap(buf)])
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		buf = buf[:n]

		p.parseFullDocument(buf)
	}

	return nil
}
