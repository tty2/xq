/*
Package attributes is responsible for parsing and printing attributes data.
*/
package attributes

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/tty2/xq/internal/domain"
)

type (
	// Processor is an attrubute processor. Keeps needed attributes to process data and handle attribute data.
	Processor struct {
		queryPath []domain.Step
		attribute string
		printList []string
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
		}
	}()

	return nil
}

func (p *Processor) process(chunk []byte) error {
	return nil
}
