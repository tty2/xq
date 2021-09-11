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
		insideTag      bool
		queryPath      []domain.Step
		currentPath    []string
		currentTag     tag
		targetTagsList []string
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

		err = p.process(buf)
		if err != nil {
			return err
		}
	}

	p.printTagsInside()

	return nil
}

func (p *Processor) printTagsInside() {
	for i := range p.targetTagsList {
		fmt.Println(p.targetTagsList[i]) // nolint forbidigo: the purpose of the function is print to stdout
	}
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

	p.updateTagList()

	if p.currentTag.skip {
		return nil
	}

	return p.updatePath()
}

func (p *Processor) processCurrentTag() error {
	if len(p.currentTag.bytes) < 3 {
		return errors.New("tag can't be less then 3 bytes")
	}

	if p.currentTag.bytes[0] != symbol.OpenBracket {
		return errors.New("tag must start from open bracket symbol")
	}

	p.markIfSkip()

	startName := 1                    // name starts after open bracket
	if p.currentTag.bytes[1] == '/' { // closed tag
		startName = 2
		p.currentTag.closed = true
	}

	endName := startName

	for ; endName < len(p.currentTag.bytes)-1; endName++ {
		if p.currentTag.bytes[endName] == ' ' {
			break
		}
	}

	p.currentTag.name = string(p.currentTag.bytes[startName:endName])

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

func (p *Processor) updateTagList() {
	if !domain.PathsMatch(p.queryPath, p.currentPath) {
		return
	}

	if slice.ContainsString(p.targetTagsList, p.currentTag.name) {
		return
	}
	// step back after deeper nesting tag with close tag (queryPath == currentPath)
	if p.queryPath[len(p.queryPath)-1].Name == p.currentTag.name {
		return
	}
	p.targetTagsList = append(p.targetTagsList, p.currentTag.name)
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
