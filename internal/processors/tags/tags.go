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
	"github.com/tty2/xq/pkg/slice"
)

const (
	closeBracket = '>'
	openBracket  = '<'
)

type (
	Processor struct {
		insideTag      bool
		targetTagFound bool
		index          int
		path           []domain.Step
		currentTag     tag
		targetTagsList []string
	}

	tag struct {
		bytes []byte
	}
)

func NewProcessor(path []domain.Step) *Processor {
	return &Processor{
		path: path,
	}
}

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

	p.printTagsInside()

	return nil
}

func (p *Processor) printTagsInside() {
	for i := range p.targetTagsList {
		fmt.Println(p.targetTagsList[i])
	}
}

func (p *Processor) process(chunk []byte) error {
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

				if tagName != p.path[p.index].Name {
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