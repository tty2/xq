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

var errServiceTag = errors.New("not a tag or service tag")

type (
	// Processor is a tag processor. Keeps needed attributes to process data and handle tag data.
	Processor struct {
		insideTag    bool
		currentPath  []string
		printList    []string
		printtedList []string
		currentTag   tag
		query        query
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
		if errors.Is(err, errServiceTag) {
			return nil
		}

		return err
	}

	p.updatePrintList()

	if p.currentTag.bytes[len(p.currentTag.bytes)-2] == '/' { // singe tag
		return nil
	}

	return p.updatePath()
}

func (p *Processor) processCurrentTag() error {
	// cdata or service tags
	if p.currentTag.bytes[1] == '?' || p.currentTag.bytes[1] == '!' {
		return errServiceTag
	}

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

func (p *Processor) updatePrintList() {
	if !domain.PathsMatch(p.query.path, p.currentPath) {
		return
	}

	if slice.ContainsString(p.printtedList, p.currentTag.name) {
		return
	}
	// step back after deeper nesting tag with close tag (queryPath == currentPath)
	if p.query.path[len(p.query.path)-1].Name == p.currentTag.name {
		return
	}

	if p.query.searchType == domain.TagList {
		p.printList = append(p.printList, p.currentTag.name)
		p.printtedList = append(p.printtedList, p.currentTag.name)
	} else if p.query.searchType == domain.AttrList {
		p.printList = pickAttributesNames(p.currentTag.bytes)
	}
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

func pickAttributesNames(tag []byte) []string {
	if len(tag) == 0 || tag[0] != symbol.OpenBracket {
		return nil
	}

	var i int
	// skip tag name
	for ; i < len(tag) && tag[i] != ' '; i++ {
	}

	var isAttrName bool
	attrs := []string{}
	attrName := []byte{}
	quotes := []byte{}
	for ; i < len(tag); i++ {
		if symbol.IsQuote(tag[i]) {
			if len(quotes) == 0 {
				quotes = append(quotes, tag[i])
				continue
			}

			if tag[i] == quotes[0] {
				quotes = []byte{}
			} else {
				quotes = append(quotes, tag[i])
			}
		}

		if len(quotes) > 0 {
			continue
		}

		if tag[i] == '=' {
			attrs = append(attrs, string(attrName))
			isAttrName = false
			attrName = []byte{}
			continue
		}
		if isAttrName {
			attrName = append(attrName, tag[i])
			continue
		}
		if tag[i] == ' ' {
			isAttrName = true
		}
	}

	return attrs
}
