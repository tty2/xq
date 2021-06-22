package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	closeBracket = '>'
	openBracket  = '<'

	minTagSize = 3 // minimum tag size can be 3. as example <b>

	green = "\033[01;32m"
	white = "\033[00m"

	maxLineLen     = 120
	indentItemSize = 2

	carriageReturn = 10
)

type parser struct {
	CurrentTag     tag
	Data           []byte
	IndentItemSize int
	Indentation    int
	InsideTag      bool // semaphore that shows if we read data inside a tag
	SkipData       bool
}

type tag struct {
	Name     string
	String   string
	Bytes    []byte
	Brackets int
}

func NewParser() parser {
	return parser{
		IndentItemSize: indentItemSize,
		Data:           []byte{},
		CurrentTag: tag{
			Bytes: []byte{},
		},
	}
}

func (p *parser) process(chunk []byte) {

	for i := range chunk {
		// skip carriage return and new line from data in order do not duplicate with created ones by parser
		if p.SkipData && (chunk[i] == ' ' || chunk[i] == '\t') {
			continue
		}
		if chunk[i] == carriageReturn {
			p.SkipData = true
			continue
		}
		p.SkipData = false

		if p.InsideTag {
			p.CurrentTag.Bytes = append(p.CurrentTag.Bytes, chunk[i])

			if chunk[i] == closeBracket {
				p.CurrentTag.Brackets -= 1

				if p.CurrentTag.Brackets > 0 {
					continue
				}

				p.InsideTag = false
				p.Data = []byte{}
				p.printTag()
			} else if chunk[i] == openBracket {
				p.CurrentTag.Brackets += 1
			}

			continue
		}

		if chunk[i] == openBracket {

			p.InsideTag = true
			p.CurrentTag = tag{
				Bytes: []byte{chunk[i]},
			}
			p.CurrentTag.Brackets += 1

			if len(p.Data) > 0 {
				fmt.Printf("\n%s", strings.Repeat("  ", p.Indentation)+string(p.Data))
			}

			continue
		}

		p.Data = append(p.Data, chunk[i])
	}
}

func (p *parser) printTag() {
	if len(p.CurrentTag.Bytes) < minTagSize {
		log.Fatalf("tag size is too small = %d, tag is `%s`", len(p.CurrentTag.Bytes), p.CurrentTag.Bytes)
	}
	if p.CurrentTag.Bytes[1] == '!' || p.CurrentTag.Bytes[1] == '?' {
		fmt.Printf("%s", p.CurrentTag.Bytes)
		return
	}

	startName := 1
	if p.CurrentTag.Bytes[1] == '/' {
		startName = 2
		p.Indentation--
		defer p.downIndent()
	}

	endName := startName
	for ; endName < len(p.CurrentTag.Bytes)-1; endName++ {
		if p.CurrentTag.Bytes[endName] == ' ' {
			break
		}
	}

	p.CurrentTag.String = strings.Repeat("  ", p.Indentation) + string(p.CurrentTag.Bytes[:startName]) +
		green + string(p.CurrentTag.Bytes[startName:endName]) +
		white + string(p.CurrentTag.Bytes[endName:len(p.CurrentTag.Bytes)])

	fmt.Printf("\n%s", p.CurrentTag.String)

	p.Indentation++
}

func (p *parser) downIndent() {
	p.Indentation--
}
