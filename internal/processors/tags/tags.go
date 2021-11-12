/*
Package tags is responsible for parsing and printing tags data.
*/
package tags

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tty2/xq/internal/domain"
	"github.com/tty2/xq/internal/domain/symbol"
	"github.com/tty2/xq/pkg/slice"
)

const indentItemSize int = 2

type (
	// Processor is a tag processor. Keeps needed attributes to process data and handle tag data.
	Processor struct {
		insideTag   bool
		currentPath []string
		printList   []string
		currentTag  tag
		query       query
		indentation int
		tagValue    []byte
		stop        bool
		index       index
	}

	query struct {
		path       []domain.Step
		attribute  string
		searchType domain.SearchType
	}

	index struct {
		set          bool
		insideTarget bool
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
		index: index{
			set: isIndexSearch(path),
		},
	}, nil
}

func isIndexSearch(path []domain.Step) bool {
	for i := range path {
		if path[i].Index > -1 {
			return true
		}
	}

	return false
}

func (p *Processor) indexTagFound() bool {
	if len(p.query.path) != len(p.currentPath) {
		return false
	}

	if !domain.PathsMatch(p.query.path, p.currentPath) {
		return false
	}

	for i := range p.query.path {
		if p.query.path[i].Index > -1 {
			return false
		}
	}

	if !p.index.insideTarget {
		p.index.insideTarget = true
	}

	return true
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

			if p.stop {
				return
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
		switch {
		case p.insideTag:
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
			if p.stop {
				return nil
			}
		case chunk[i] == symbol.OpenBracket:
			p.insideTag = true
			p.currentTag = tag{
				bytes:    []byte{symbol.OpenBracket},
				brackets: 1,
			}
			if p.query.searchType == domain.TagValue && p.intoQueryPath() {
				if p.index.set {
					if !p.index.insideTarget {
						continue
					}
				}
				if strings.TrimSpace(string(p.tagValue)) == "" {
					p.tagValue = []byte{}

					continue
				}
				p.printList = append(p.printList, string(append(bytes.Repeat([]byte(" "),
					indentItemSize*p.indentation+indentItemSize), p.tagValue...)))
				p.tagValue = []byte{}
			}
		case p.query.searchType == domain.TagValue && p.intoQueryPath():
			if p.index.set {
				if !p.index.insideTarget {
					continue
				}
			}
			if chunk[i] == symbol.NewLine || chunk[i] == symbol.CarriageReturn {
				continue
			}
			p.tagValue = append(p.tagValue, chunk[i])
		}
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
		if p.index.set && // process with index in path initialized
			!p.stop && // target isn't parsed completely
			p.index.insideTarget && // but current tag is target
			domain.PathsMatch(p.query.path, p.currentPath) { // and this is the close tag for target
			p.updatePrintList()
			p.stop = true
		}
	} else {
		p.currentPath = append(p.currentPath, p.currentTag.name)
	}

	p.updatePrintList()
	if p.currentTagIsSingle() {
		p.currentPath = p.currentPath[:len(p.currentPath)-1]
	}
	if p.currentTag.closed {
		return p.decrementPath()
	}

	return nil
}

func (p *Processor) updatePrintList() {
	if p.index.set {
		if p.stop {
			return
		}

		if !p.currentTag.closed && !p.index.insideTarget {
			p.decrementSearchIndex()
			p.indexTagFound()
		}

		if !p.index.insideTarget {
			return
		}
	}
	switch {
	case p.query.searchType == domain.TagList && p.tagInQueryPath():
		tn := strings.TrimSpace(p.currentTag.name)
		if !slice.ContainsString(p.printList, tn) {
			p.printList = append(p.printList, tn)
		}
	case p.query.searchType == domain.AttrList && domain.PathsMatch(p.query.path, p.currentPath):
		list := pickAttributesNames(p.currentTag.bytes)
		for i := range list {
			an := strings.TrimSpace(list[i])
			if !slice.ContainsString(p.printList, an) {
				p.printList = append(p.printList, an)
			}
		}
	case p.query.searchType == domain.AttrValue && domain.PathsMatch(p.query.path, p.currentPath):
		av, err := pickAttributeValue(p.query.attribute, p.currentTag.bytes)
		if err != nil {
			return
		}
		if av == "" {
			return
		}
		if !slice.ContainsString(p.printList, av) {
			p.printList = append(p.printList, av)
		}
	case p.query.searchType == domain.TagValue && p.intoQueryPath():
		p.indentation = len(p.currentPath) - len(p.query.path)
		p.printList = append(p.printList, string(append(bytes.Repeat([]byte(" "),
			indentItemSize*p.indentation), p.currentTag.bytes...)))
	}
}

func (p *Processor) tagInQueryPath() bool {
	// +1 because /query/path/tag + current_tag
	if len(p.query.path)+1 != len(p.currentPath) {
		return false
	}

	// -1 because len - current tag
	for i := 0; i < len(p.currentPath)-1; i++ {
		if p.query.path[i].Name != p.currentPath[i] {
			return false
		}
	}

	return true
}

func (p *Processor) intoQueryPath() bool {
	if len(p.query.path) > len(p.currentPath) {
		return false
	}

	for i := 0; i < len(p.query.path); i++ {
		if p.query.path[i].Name != p.currentPath[i] {
			return false
		}
	}

	return true
}

func (p *Processor) pathIntoQuery() bool {
	if len(p.currentPath) > len(p.query.path) {
		return false
	}

	for i := 0; i < len(p.currentPath); i++ {
		if p.query.path[i].Name != p.currentPath[i] {
			return false
		}
	}

	return true
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

func (p *Processor) decrementSearchIndex() {
	if len(p.currentPath) > len(p.query.path) {
		return
	}
	if len(p.currentPath) == 0 {
		return
	}
	if p.query.path[len(p.currentPath)-1].Index == -1 {
		return
	}

	if p.pathIntoQuery() {
		p.query.path[len(p.currentPath)-1].Index--
	}
}
