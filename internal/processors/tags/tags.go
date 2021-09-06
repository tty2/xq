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
		bytes  []byte
		name   string
		closed bool
		single bool
	}
)

// NewProcessor creates a new Processor with needed attributes.
func NewProcessor(path []domain.Step) *Processor {
	return &Processor{
		queryPath: path,
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
			p.currentTag.bytes = append(p.currentTag.bytes, chunk[i])

			if chunk[i] != symbol.CloseBracket {
				continue
			}
			p.insideTag = false

			err := p.processCurrentTag()
			if err != nil {
				return err
			}

			err = p.updatePath()
			if err != nil {
				return err
			}

			p.updateTagList()
		}

		if chunk[i] == symbol.OpenBracket {
			p.insideTag = true
			p.currentTag = tag{
				bytes: []byte{chunk[i]},
			}
		}
	}

	return nil
}

func (p *Processor) processCurrentTag() error {
	if len(p.currentTag.bytes) < 3 {
		return errors.New("tag can't be less then 3 bytes")
	}

	if p.currentTag.bytes[0] != symbol.OpenBracket {
		return errors.New("tag must start from open bracket symbol")
	}

	startName := 1                    // name starts after open bracket
	if p.currentTag.bytes[1] == '/' { // closed tag
		startName = 2
		p.currentTag.closed = true
	}
	if p.currentTag.bytes[len(p.currentTag.bytes)-2] == '/' {
		p.currentTag.single = true
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

func (p *Processor) updateTagList() {
	if !domain.PathsMatch(p.queryPath, p.currentPath) {
		return
	}

	if !slice.ContainsString(p.targetTagsList, p.currentTag.name) {
		return
	}

	p.targetTagsList = append(p.targetTagsList, p.currentTag.name)
}

func (p *Processor) updatePath() error {
	if p.currentTag.single {
		return nil
	}

	lastElement := len(p.currentPath) - 1

	if p.currentTag.closed {
		if p.currentPath[lastElement] != p.currentTag.name {
			return fmt.Errorf("incorrect xml structure: the last open tag is %s, but close tag is %s",
				p.currentPath[lastElement], p.currentTag.name)
		}

		p.currentPath = p.currentPath[:lastElement]

		return nil
	}

	p.currentPath = append(p.currentPath, p.currentTag.name)

	return nil
}
