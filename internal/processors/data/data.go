/*
Package data is responsible for full document parse.
This part of processing is responsible for full document parsing with colorizing and
making indentation. This processing doesn't parse and separate tags, attributes and
values from document structure but optimized for fast processing with less memory
consumption.
*/
package data

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/tty2/xq/internal/domain/color"
	"github.com/tty2/xq/internal/domain/symbol"
)

const (
	minTagSize = 3 // minimum tag size can be 3. as example <b>
)

type (
	// Processor is a data processor. Keeps needed attributes to process data, colorize and print it.
	Processor struct {
		CurrentTag     tag
		Data           []byte
		IndentItemSize int
		Indentation    int
		InsideTag      bool // semaphore that shows if we read data inside a tag
		SkipData       bool
	}

	tag struct {
		Name     string
		String   string
		Bytes    []byte
		Brackets int
	}

	attribute struct {
		Name        []byte
		Value       []byte
		Quote       byte
		NextIsQuote bool
		InsideValue bool
	}
)

// NewProcessor creates a new Processor with needed attributes.
func NewProcessor(indentationSize int) *Processor {
	return &Processor{
		IndentItemSize: indentationSize,
	}
}

// Process reads the data from `r` reader and processes it.
func (p *Processor) Process(r *bufio.Reader) error {
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

		p.process(buf)
	}

	return nil
}

func (p *Processor) process(chunk []byte) {
	for i := range chunk {
		// skip carriage return and new line from data in order do not duplicate with created ones by Processor
		if p.SkipData && (chunk[i] == ' ' || chunk[i] == '\t') {
			continue
		}
		if chunk[i] == symbol.NewLine || chunk[i] == symbol.CarriageReturn {
			p.SkipData = true

			continue
		}
		p.SkipData = false

		if p.InsideTag {
			p.CurrentTag.Bytes = append(p.CurrentTag.Bytes, chunk[i])

			if chunk[i] == symbol.CloseBracket {
				p.CurrentTag.Brackets--

				if p.CurrentTag.Brackets > 0 {
					continue
				}

				p.InsideTag = false
				p.printTag()
				p.SkipData = true // skip if there are empty symbol beeween close tag and new data
			} else if chunk[i] == symbol.OpenBracket {
				p.CurrentTag.Brackets++
			}

			continue
		}

		if chunk[i] == symbol.OpenBracket {
			p.InsideTag = true
			p.CurrentTag = tag{
				Bytes: []byte{chunk[i]},
			}
			p.CurrentTag.Brackets++

			if len(p.Data) > 0 {
				// nolint forbidigo: printf is executed on purpose here
				fmt.Printf("%s\n", append(bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation), p.Data...))
				p.Data = []byte{}
			}

			continue
		}

		p.Data = append(p.Data, chunk[i])
	}
}

// nolint forbidigo: printf in this method is executed on purpose
func (p *Processor) printTag() {
	if len(p.CurrentTag.Bytes) < minTagSize {
		log.Fatalf("tag size is too small = %d, tag is `%s`", len(p.CurrentTag.Bytes), p.CurrentTag.Bytes)
	}
	if p.CurrentTag.Bytes[1] == '!' || p.CurrentTag.Bytes[1] == '?' { // service tag, comment or cdata
		fmt.Printf("%s\n", append(bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation), p.CurrentTag.Bytes...))

		return
	}

	fmt.Printf("%s\n", p.colorizeTag())

	if p.CurrentTag.Bytes[len(p.CurrentTag.Bytes)-2] != '/' {
		p.Indentation++
	}
}

func (p *Processor) downIndent() {
	p.Indentation--
}

func (p *Processor) colorizeTag() []byte {
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
	coloredTag = append(coloredTag, []byte(color.Red)...)                                         // add red color
	coloredTag = append(coloredTag, p.CurrentTag.Bytes[startName:endName]...)                     // tag name

	attr := attribute{
		Value: []byte{},
	}
	for i := endName; i < ln-1; i++ {
		if attr.NextIsQuote {
			if symbol.IsQuote(p.CurrentTag.Bytes[i]) {
				attr.Quote = p.CurrentTag.Bytes[i]
				attr.NextIsQuote = false
				attr.InsideValue = true
			}

			continue
		}
		if attr.InsideValue {
			if p.CurrentTag.Bytes[i] == attr.Quote && attr.Value[len(attr.Value)-1] != '\\' {
				attr.InsideValue = false
				coloredTag = append(coloredTag, symbol.Space)
				coloredTag = append(coloredTag, []byte(color.Green)...)
				coloredTag = append(coloredTag, attr.Name...)
				coloredTag = append(coloredTag, []byte(color.White)...)
				coloredTag = append(coloredTag, '=', attr.Quote)
				coloredTag = append(coloredTag, attr.Value...)
				coloredTag = append(coloredTag, attr.Quote)
				attr = attribute{
					Value: []byte{},
				}
			} else {
				attr.Value = append(attr.Value, p.CurrentTag.Bytes[i])
			}

			continue
		}
		if p.CurrentTag.Bytes[i] == ' ' {
			continue
		}
		if p.CurrentTag.Bytes[i] == '=' { // value of attribute
			attr.NextIsQuote = true
			coloredTag = append(coloredTag, []byte(color.White)...)

			continue
			// i != ln-3 in order do not colorize `/` sign inside an empty tag in case like this `<...attr="value" />`
		} else if p.CurrentTag.Bytes[i] == '/' && i == ln-2 { // end attribute value
			coloredTag = append(coloredTag, p.CurrentTag.Bytes[i])

			continue
		}
		attr.Name = append(attr.Name, p.CurrentTag.Bytes[i])
	}
	coloredTag = append(coloredTag, []byte(color.White)...)
	coloredTag = append(coloredTag, p.CurrentTag.Bytes[ln-1]) // add close bracket

	return coloredTag
}
