package attributes

import (
	"bufio"
	"fmt"
	"io"
)

type (
	Processor struct {
		targetAttributesList []string
		path                 []string
		attribute            string
	}
)

func NewProcessor(path []string, attribute string) *Processor {
	return &Processor{
		path:      path,
		attribute: attribute,
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

	p.printAttrubutes()

	return nil
}

func (p *Processor) printAttrubutes() {
	for i := range p.targetAttributesList {
		fmt.Println(p.targetAttributesList[i])
	}
}

func (p *Processor) process(chunk []byte) {}
