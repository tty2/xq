/*
Package attributes is responsible for parsing and printing attributes data.
*/
package attributes

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tty2/xq/internal/domain"
)

type (
	// Processor is an attrubute processor. Keeps needed attributes to process data and handle attribute data.
	Processor struct {
		queryPath            []domain.Step
		currentPath          []string
		attribute            string
		targetAttributesList []string
	}
)

// NewProcessor creates a new Processor with needed attributes.
func NewProcessor(path []domain.Step, attribute string) (*Processor, error) {
	if len(path) == 0 {
		return nil, errors.New("query path must not be empty")
	}
	if strings.TrimSpace(attribute) == "" {
		return nil, errors.New("attribute you search must not be empty")
	}

	return &Processor{
		queryPath: path,
		attribute: attribute,
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

		p.process(buf)
	}

	p.printAttrubutes()

	return nil
}

func (p *Processor) printAttrubutes() {
	for i := range p.targetAttributesList {
		fmt.Println(p.targetAttributesList[i]) // nolint forbidigo: the purpose of the function is print to stdout
	}
}

func (p *Processor) process(chunk []byte) {

}
