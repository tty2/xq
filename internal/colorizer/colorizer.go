/*
Package datcolorizera is responsible for full document parse.
This part of processing is responsible for full document parsing with colorizing and
making indentation. This processing doesn't parse and separate tags, attributes and
values from document structure but optimized for fast processing with less memory
consumption.
*/
package colorizer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/tty2/xq/internal/domain"
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
	if p.CurrentTag.Bytes[1] == '/' {
		p.Indentation--
		defer p.downIndent()
	}

	coloredTag := domain.ColorizeTag(p.CurrentTag.Bytes)
	coloredTag = append(bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation), coloredTag...) // add indentation
	p.printList = append(p.printList, string(coloredTag))

	if p.CurrentTag.Bytes[len(p.CurrentTag.Bytes)-2] != '/' {
		p.Indentation++
	}

	return nil
}

func (p *Processor) downIndent() {
	p.Indentation--
}
