/*
Package tags is responsible for parsing and printing tags data.
*/
package tags

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/tty2/xq/internal/domain"
	"github.com/tty2/xq/internal/domain/symbol"
	"github.com/tty2/xq/pkg/slice"
)

type (
	// Processor is a tag processor. Keeps needed attributes to process data and handle tag data.
	Processor struct {
		insideTag bool
		queryPath []domain.Step
		// queryAttribute   string
		currentPath      []string
		currentTag       tag
		printtedTagsList []string
		printList        []string
	}

	tag struct {
		bytes    []byte
		name     string
		closed   bool
		skip     bool
		brackets int
	}
)

// NewProcessor creates a new Processor with needed attributes.
func NewProcessor(path []domain.Step) (*Processor, error) {
	if len(path) == 0 {
		return nil, errors.New("query path must not be empty")
	}

	return &Processor{
		queryPath: path,
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
					break
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
		if p.insideTag {
			err := p.addSymbolIntoTag(chunk[i])
			if err != nil {
				return err
			}
		} else if chunk[i] == symbol.OpenBracket {
			p.insideTag = true
			p.currentTag = tag{
				bytes:    []byte{symbol.OpenBracket},
				brackets: 1,
			}
		}
		// skip data outside tags because we are interested in tags only
	}

	return nil
}

func (p *Processor) addSymbolIntoTag(s byte) error {
	if s == symbol.OpenBracket {
		p.currentTag.brackets++

		return nil
	}

	p.currentTag.bytes = append(p.currentTag.bytes, s)

	if s != symbol.CloseBracket {
		return nil
	}

	p.currentTag.brackets--
	if p.currentTag.brackets > 0 {
		return nil
	}

	p.insideTag = false

	err := p.processCurrentTag()
	if err != nil {
		return err
	}

	p.updatePrintList()

	if p.currentTag.skip {
		return nil
	}

	return p.updatePath()
}

func (p *Processor) processCurrentTag() error {
	p.markIfSkip()

	tg := domain.Tag{
		Bytes: p.currentTag.bytes,
	}

	err := tg.SetName()
	if err != nil {
		return err
	}

	p.currentTag.closed = p.currentTag.bytes[1] == '/'
	p.currentTag.name = tg.Name

	return nil
}

func (p *Processor) markIfSkip() {
	if p.currentTag.bytes[len(p.currentTag.bytes)-2] == '/' {
		p.currentTag.skip = true
	}
	if p.currentTag.bytes[1] == '?' {
		p.currentTag.skip = true
	}
	if p.currentTag.bytes[1] == '!' {
		p.currentTag.skip = true
	}
}

func (p *Processor) updatePrintList() {
	if !domain.PathsMatch(p.queryPath, p.currentPath) {
		return
	}

	if slice.ContainsString(p.printtedTagsList, p.currentTag.name) {
		return
	}
	// step back after deeper nesting tag with close tag (queryPath == currentPath)
	if p.queryPath[len(p.queryPath)-1].Name == p.currentTag.name {
		return
	}
	p.printList = append(p.printList, p.currentTag.name)
	p.printtedTagsList = append(p.printtedTagsList, p.currentTag.name)
}

func (p *Processor) updatePath() error {
	lastElement := len(p.currentPath) - 1

	if p.currentTag.closed {
		if p.currentPath[lastElement] != p.currentTag.name {
			return fmt.Errorf("incorrect xml structure: the last open tag is %s, but close tag is %s",
				p.currentPath[lastElement], p.currentTag.name)
		}

		p.currentPath = p.currentPath[:lastElement] // TODO: consider to add p.currentTag.brackets-- after this line

		return nil
	}

	p.currentPath = append(p.currentPath, p.currentTag.name)

	return nil
}
