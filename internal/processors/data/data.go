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
		printList      []string
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
func NewProcessor(indentationSize int) (*Processor, error) {
	return &Processor{
		IndentItemSize: indentationSize,
	}, nil
}

// Process reads the data from `r` reader and processes it.
func (p *Processor) Process(r *bufio.Reader) chan string {
	buf := make([]byte, 0, 4*1024)
	ch := make(chan string)

	go func() {
		defer close(ch)
		for {
			n, err := r.Read(buf[:cap(buf)])
			if err != nil {
				if err == io.EOF {
					return
				}
				ch <- err.Error()

				return
			}

			buf = buf[:n]

			err = p.process(buf)
			if err != nil {
				ch <- err.Error()

				return
			}

			for i := range p.printList {
				ch <- p.printList[i]
			}

			p.printList = []string{}
		}
	}()

	return ch
}

func (p *Processor) process(chunk []byte) error {
	for i := range chunk {
		switch {
		case p.skip(chunk[i]):
			continue
		case p.InsideTag:
			err := p.closeTag(chunk[i])
			if err != nil {
				return err
			}
		case chunk[i] == symbol.OpenBracket:
			p.startTag(chunk[i])
		default:
			p.Data = append(p.Data, chunk[i])
		}
	}

	return nil
}

func (p *Processor) skip(b byte) bool {
	// skip carriage return and new line from data in order do not duplicate with created ones by Processor
	if p.SkipData && (b == ' ' || b == '\t') {
		return true
	}
	if b == symbol.NewLine || b == symbol.CarriageReturn {
		p.SkipData = true

		return true
	}

	p.SkipData = false

	return false
}

func (p *Processor) startTag(b byte) {
	p.InsideTag = true
	p.CurrentTag = tag{
		Bytes: []byte{b},
	}
	p.CurrentTag.Brackets++

	if len(p.Data) > 0 {
		p.printList = append(p.printList, string(append(bytes.Repeat([]byte(" "),
			p.IndentItemSize*p.Indentation), p.Data...)))
		p.Data = []byte{}
	}
}

func (p *Processor) closeTag(b byte) error {
	p.CurrentTag.Bytes = append(p.CurrentTag.Bytes, b)

	switch b {
	case symbol.CloseBracket:
		p.CurrentTag.Brackets--

		if p.CurrentTag.Brackets > 0 {
			return nil
		}

		p.InsideTag = false
		err := p.addToPrintList()
		if err != nil {
			return err
		}
		p.SkipData = true // skip if there are empty symbol beeween close tag and new data
	case symbol.OpenBracket:
		p.CurrentTag.Brackets++
	}

	return nil
}

func (p *Processor) addToPrintList() error {
	if len(p.CurrentTag.Bytes) < minTagSize {
		return fmt.Errorf("tag size is too small = %d, tag is `%s`", len(p.CurrentTag.Bytes), p.CurrentTag.Bytes)
	}
	if p.CurrentTag.Bytes[1] == '!' || p.CurrentTag.Bytes[1] == '?' { // service tag, comment or cdata
		p.printList = append(p.printList, string(
			append(bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation), p.CurrentTag.Bytes...)))

		return nil
	}

	p.printList = append(p.printList, string(p.colorizeTag()))

	if p.CurrentTag.Bytes[len(p.CurrentTag.Bytes)-2] != '/' {
		p.Indentation++
	}

	return nil
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
			if p.CurrentTag.Bytes[i] == attr.Quote && (len(attr.Value) == 0 || attr.Value[len(attr.Value)-1] != '\\') {
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
