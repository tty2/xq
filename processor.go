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

	red   = "\033[01;31m"
	green = "\033[01;32m"
	white = "\033[00m"

	indentItemSize = 2

	carriageReturn = 10 // '\n'
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

func newParser() parser {
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
				p.CurrentTag.Brackets--

				if p.CurrentTag.Brackets > 0 {
					continue
				}

				p.InsideTag = false
				p.Data = []byte{}
				p.printTag()
			} else if chunk[i] == openBracket {
				p.CurrentTag.Brackets++
			}

			continue
		}

		if chunk[i] == openBracket {

			p.InsideTag = true
			p.CurrentTag = tag{
				Bytes: []byte{chunk[i]},
			}
			p.CurrentTag.Brackets++

			if len(p.Data) > 0 {
				fmt.Printf("\n%s", strings.Repeat("  ", p.Indentation)+string(p.Data))
			}

			continue
		}

		p.Data = append(p.Data, chunk[i])
	}
}

func (p *parser) printTag() {
	ln := len(p.CurrentTag.Bytes)
	if ln < minTagSize {
		log.Fatalf("tag size is too small = %d, tag is `%s`", ln, p.CurrentTag.Bytes)
	}
	if p.CurrentTag.Bytes[1] == '!' || p.CurrentTag.Bytes[1] == '?' { // service tag or cdata
		fmt.Printf("%s", p.CurrentTag.Bytes)
		return
	}

	printBytes := make([]byte, 0, p.Indentation+ln)

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

	for l := p.IndentItemSize * p.Indentation; l > 0; l-- {
		printBytes = append(printBytes, ' ')
	}

	printBytes = append(printBytes, p.CurrentTag.Bytes[:startName]...)        // add indentation
	printBytes = append(printBytes, []byte(red)...)                           // add red color
	printBytes = append(printBytes, p.CurrentTag.Bytes[startName:endName]...) // tag name
	printBytes = append(printBytes, []byte(green)...)

	for i := endName; i < ln-1; i++ {
		if p.CurrentTag.Bytes[i] == '=' {
			printBytes = append(printBytes, []byte(white)...)
		} else if p.CurrentTag.Bytes[i] == ' ' && p.CurrentTag.Bytes[i-1] == '"' {
			printBytes = append(printBytes, []byte(green)...)
		}
		printBytes = append(printBytes, p.CurrentTag.Bytes[i])
	}
	printBytes = append(printBytes, []byte(white)...)         // turn back to white color
	printBytes = append(printBytes, p.CurrentTag.Bytes[ln-1]) // add close bracket

	fmt.Printf("\n%s", printBytes)

	p.Indentation++
}

func (p *parser) downIndent() {
	p.Indentation--
}
