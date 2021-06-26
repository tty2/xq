package main

import (
	"bytes"
	"fmt"
	"log"
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
				fmt.Printf("%s\n", append(bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation), p.Data...))
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
	if p.CurrentTag.Bytes[1] == '!' || p.CurrentTag.Bytes[1] == '?' { // service tag, comment or cdata
		fmt.Printf("%s\n", append(bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation), p.CurrentTag.Bytes...))
		return
	}

	fmt.Printf("%s\n", p.colorizeTag())

	p.Indentation++
}

func (p *parser) downIndent() {
	p.Indentation--
}

func (p *parser) colorizeTag() []byte {
	ln := len(p.CurrentTag.Bytes)

	coloredTag := make([]byte, 0, p.Indentation+ln)

	startName := 1                    // name starts after open bracket
	if p.CurrentTag.Bytes[1] == '/' { // closed tag
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

	coloredTag = append(coloredTag, bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation)...) // add indentation
	coloredTag = append(coloredTag, p.CurrentTag.Bytes[:startName]...)                            // add open bracket
	coloredTag = append(coloredTag, []byte(red)...)                                               // add red color
	coloredTag = append(coloredTag, p.CurrentTag.Bytes[startName:endName]...)                     // tag name
	coloredTag = append(coloredTag, []byte(green)...)                                             // attribute name starts

	for i := endName; i < ln-1; i++ {
		if p.CurrentTag.Bytes[i] == '=' { // value of attribute
			coloredTag = append(coloredTag, []byte(white)...)
		} else if p.CurrentTag.Bytes[i] == ' ' && p.CurrentTag.Bytes[i-1] == '"' { // end attribute value
			coloredTag = append(coloredTag, []byte(green)...)
		}
		coloredTag = append(coloredTag, p.CurrentTag.Bytes[i])
	}
	coloredTag = append(coloredTag, []byte(white)...)
	coloredTag = append(coloredTag, p.CurrentTag.Bytes[ln-1]) // add close bracket

	return coloredTag
}
