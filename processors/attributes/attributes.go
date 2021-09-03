package attributeparser

import (
	"bufio"
	"fmt"
	"io"
)

type (
	Parser struct {
		targetAttributesList []string
	}
)

func (p *Parser) Process(r *bufio.Reader) error {
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

func (p *Parser) printAttrubutes() {
	for i := range p.targetAttributesList {
		fmt.Println(p.targetAttributesList[i])
	}
}

func (p *Parser) process(chunk []byte) {}
