/*
Package tagparser is responsible for parsing and printing tags data.
*/
package tagparser

import (
	"errors"
	"fmt"
	"io"

	"github.com/tty2/xq/pkg/slice"
)

const (
	closeBracket = '>'
	openBracket  = '<'
)

type (
	Parser struct {
		insideTag      bool
		targetTagFound bool
		index          int
		path           []string
		currentTag     tag
		targetTagsList []string
	}

	tag struct {
		bytes []byte
	}
)

func (p *Parser) Process(chunk []byte) {
	err := p.getTagsList(chunk)
	if err == io.EOF {
		p.printTagsInside()
	}
}

func (p *Parser) printTagsInside() {
	for i := range p.targetTagsList {
		fmt.Println(p.targetTagsList[i])
	}
}

func (p *Parser) getTagsList(chunk []byte) error {
	for i := range chunk {
		if p.insideTag {
			p.currentTag.bytes = append(p.currentTag.bytes, chunk[i])

			if chunk[i] == closeBracket {
				p.insideTag = false

				tagName, err := getTagName(p.currentTag.bytes)
				if err != nil {
					return err
				}

				if p.targetTagFound {
					if !slice.ContainsString(p.targetTagsList, tagName) {

						continue
					}
					p.targetTagsList = append(p.targetTagsList, tagName)

					continue
				}

				if tagName != p.path[p.index] {
					return errors.New("incorrect tag name in path")
				}

				p.index++

				if len(p.path) > p.index {
					continue
				}

				p.targetTagFound = true
			}

			continue
		}

		if chunk[i] == openBracket {
			p.insideTag = true
			p.currentTag = tag{
				bytes: []byte{chunk[i]},
			}

			continue
		}
	}

	return nil
}

func getTagName(t []byte) (string, error) {
	if len(t) < 3 {
		return "", errors.New("tag can't be less then 3 bytes")
	}

	if t[0] != openBracket {
		return "", errors.New("tag must start from open bracket symbol")
	}

	startName := 1   // name starts after open bracket
	if t[1] == '/' { // closed tag
		startName = 2
	}

	endName := startName

	for ; endName < len(t)-1; endName++ {
		if t[endName] == ' ' {
			break
		}
	}

	return string(t[startName:endName]), nil
}
