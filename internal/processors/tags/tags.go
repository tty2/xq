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
		insideTag   bool
		currentPath []string
		printList   []string
		currentTag  tag
		query       query
	}

	query struct {
		path       []domain.Step
		attribute  string
		searchType domain.SearchType
	}

	tag struct {
		bytes    []byte
		name     string // tagname
		closed   bool   // </tagname>
		brackets int    // stack to keep track open and close brackets
	}
)

// NewProcessor creates a new Processor with needed attributes.
func NewProcessor(path []domain.Step, attribute string, search domain.SearchType) (*Processor, error) {
	if len(path) == 0 {
		return nil, errors.New("query path must not be empty")
	}

	return &Processor{
		query: query{
			path:       path,
			attribute:  attribute,
			searchType: search,
		},
	}, nil
}

// Process reads the data from `r` reader and processes it.
func (p *Processor) Process(r *bufio.Reader) chan string {
	buf := make([]byte, 0, 4*1024)
	ch := make(chan string)

	var idx int
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

			for ; idx < len(p.printList); idx++ {
				ch <- p.printList[idx]
			}
		}

		for ; idx < len(p.printList); idx++ {
			ch <- p.printList[idx]
		}
	}()

	return ch
}

func (p *Processor) process(chunk []byte) error {
	for i := range chunk {
		if p.insideTag {
			p.addSymbolIntoTag(chunk[i])

			if p.currentTag.brackets > 0 {
				continue
			}
			p.insideTag = false

			if p.skip() {
				continue
			}

			err := p.processCurrentTag()
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

func (p *Processor) addSymbolIntoTag(s byte) {
	if s == symbol.OpenBracket {
		p.currentTag.brackets++
	}

	p.currentTag.bytes = append(p.currentTag.bytes, s)

	if s != symbol.CloseBracket {
		return
	}
	p.currentTag.brackets--
}

func (p *Processor) skip() bool {
	return p.currentTag.bytes[1] == '?' ||
		p.currentTag.bytes[1] == '!'
}

func (p *Processor) currentTagIsSingle() bool {
	ln := len(p.currentTag.bytes)
	return ln > 3 && p.currentTag.bytes[ln-2] == '/'
}

func (p *Processor) processCurrentTag() error {
	err := p.currentTag.setName()
	if err != nil {
		return err
	}

	p.currentTag.closed = p.currentTag.bytes[1] == '/'

	if p.currentTag.closed {
		return p.updatePath()
	}

	if p.currentTagIsSingle() {
		p.currentPath = append(p.currentPath, p.currentTag.name)
		p.updatePrintList()
		p.currentPath = p.currentPath[:len(p.currentPath)-1]
		return nil
	}

	err = p.updatePath()
	if err != nil {
		return err
	}

	p.updatePrintList()

	return nil
}

func (p *Processor) updatePrintList() {
	if p.query.searchType == domain.TagList &&
		p.tagInQueryPath() &&
		!slice.ContainsString(p.printList, p.currentTag.name) {
		p.printList = append(p.printList, p.currentTag.name)
	} else if p.query.searchType == domain.AttrList &&
		domain.PathsMatch(p.query.path, p.currentPath) {
		list := pickAttributesNames(p.currentTag.bytes)
		for i := range list {
			if !slice.ContainsString(p.printList, list[i]) {
				p.printList = append(p.printList, list[i])
			}
		}
	}
}

func (p *Processor) tagInQueryPath() bool {
	// +1 because /query/path/tag + current_tag
	if len(p.query.path)+1 != len(p.currentPath) {
		return false
	}

	// -2 because len - current tag
	for i := 0; i < len(p.currentPath)-1; i++ {
		if p.query.path[i].Name != p.currentPath[i] {
			return false
		}
	}

	return true
}

func (p *Processor) updatePath() error {
	if p.currentTag.closed {
		return p.decrementPath()
	}

	p.currentPath = append(p.currentPath, p.currentTag.name)

	return nil
}

func (p *Processor) decrementPath() error {
	ln := len(p.currentPath)
	if ln == 0 {
		return nil
	}

	if p.currentPath[ln-1] != p.currentTag.name {
		return fmt.Errorf("incorrect xml structure: the last open tag is %s, but close tag is %s",
			p.currentPath[ln-1], p.currentTag.name)
	}

	p.currentPath = p.currentPath[:ln-1]

	return nil
}

func (t *tag) setName() error {
	if len(t.bytes) < 3 {
		return domain.ErrTagShort
	}

	if t.bytes[0] != symbol.OpenBracket {
		return domain.ErrTagShort
	}

	if t.bytes[len(t.bytes)-1] != symbol.CloseBracket {
		return domain.ErrTagShort
	}

	startName := 1         // name starts after open bracket
	if t.bytes[1] == '/' { // closed tag
		startName = 2
	}

	endName := startName

	for ; endName < len(t.bytes)-1; endName++ {
		if t.bytes[endName] == ' ' {
			break
		}
	}

	t.name = string(t.bytes[startName:endName])

	return nil
}
