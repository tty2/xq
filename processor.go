package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	closeBracket = '>'
	openBracket  = '<'

	green = "\033[01;32m"
	white = "\033[00m"

	maxLineLen     = 120
	indentItemSize = 2

	carriageReturn = '\n'
)

type parser struct {
	CurrentTag     tag
	Data           []byte
	IndentItemSize int // 2
	Indentation    int
	InsideTag      bool // semaphore that shows if we read data inside another tag
	MaxLen         int  // maximum line len
}

type tag struct {
	Name   string
	String string
	Bytes  []byte
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
		if chunk[i] == carriageReturn {
			continue
		}

		if p.InsideTag {
			p.CurrentTag.Bytes = append(p.CurrentTag.Bytes, chunk[i])

			if chunk[i] == closeBracket {
				p.InsideTag = false
				p.Data = []byte{}
				p.printTag()
				p.Indentation++
			}

			continue
		}

		if chunk[i] == openBracket {

			p.InsideTag = true
			p.CurrentTag = tag{
				Bytes: []byte{chunk[i]},
			}

			if len(p.Data) > 0 {
				fmt.Printf("\n%s", strings.Repeat("  ", p.Indentation)+string(p.Data))
			}

			continue
		}

		p.Data = append(p.Data, chunk[i])
	}

}

func (p *parser) printTag() {
	if len(p.CurrentTag.Bytes) < 3 {
		log.Fatalf("tag size is too small = %d, tag is `%s`", len(p.CurrentTag.Bytes), p.CurrentTag.Bytes)
	}
	if p.CurrentTag.Bytes[1] == '!' || p.CurrentTag.Bytes[1] == '?' {
		fmt.Printf("%s", p.CurrentTag.Bytes)
		p.Indentation--
		return
	}

	startName := 1
	if p.CurrentTag.Bytes[1] == '/' {
		startName = 2
		p.Indentation--
	}

	endName := startName
	for ; endName < len(p.CurrentTag.Bytes)-1; endName++ {
		if p.CurrentTag.Bytes[endName] == ' ' {
			break
		}
	}

	p.CurrentTag.String = strings.Repeat("  ", p.Indentation) + string(p.CurrentTag.Bytes[:startName]) + green + string(p.CurrentTag.Bytes[startName:endName]) +
		white + string(p.CurrentTag.Bytes[endName:])

	fmt.Printf("\n%s", p.CurrentTag.String)

	if p.CurrentTag.Bytes[1] == '/' {
		p.Indentation--
	}
}
