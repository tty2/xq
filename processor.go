package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	closeBracket = '>'
	openBracket  = '<'

	green = "\033[01;34m"
	white = "\033[00m"

	maxLineLen     = 120
	indentItemSize = 2
)

type parser struct {
	Tags           []string // tags stack
	CurrentTag     tag
	Indentation    int
	Interrupted    []byte // if data was started in a bytes chunk but chunk ended we need to keep the start of tag
	ReadingFrom    int    // chunk index to read from
	InsideTag      bool   // semaphore that shows if we read data inside another tag
	TagName        string // the name of tag
	Len            int    // len of line
	MaxLen         int    // maximum line len
	IndentItemSize int    // 2
	NewLine        bool
}

type tag struct {
	Name   string
	String string
}

func NewParser() parser {
	return parser{
		Tags:           []string{},
		Interrupted:    []byte{},
		IndentItemSize: indentItemSize,
	}
}

func (p *parser) process(chunk []byte) {

	ln := len(chunk)

	for i := range chunk {
		if p.InsideTag {
			if chunk[i] == closeBracket {
				p.InsideTag = false

				p.printTag(chunk[p.ReadingFrom : i+1])

				if i+1 == ln {
					p.ReadingFrom = 0
					return // chunk end is equal with tag end
				}
				p.ReadingFrom = i + 1
			}
			continue
		}

		if chunk[i] == openBracket {

			// print previous
			if p.ReadingFrom != i {
				line := fmt.Sprintf("%s%s", p.Interrupted, chunk[p.ReadingFrom:i])

				if len(line)+p.Len > p.MaxLen {
					fmt.Printf("\n%s%s", strings.Repeat("  ", p.Indentation), line)
					p.Len = 2*p.Indentation + len(line)
					p.NewLine = true
				} else {
					fmt.Printf("%s", line)
					p.Len += len(line)
				}
			}

			p.ReadingFrom = i
			p.InsideTag = true
		}
	}

	p.Interrupted = chunk[p.ReadingFrom:ln]

	p.ReadingFrom = 0
}

func (p *parser) printTag(tag []byte) {
	tg := append(p.Interrupted, tag...)
	if tg[1] == '!' {
		fmt.Printf("%s", tag)
		return
	}
	if p.NewLine {
		fmt.Printf("\n%s", strings.Repeat("  ", p.Indentation))
	}

	st := p.separateTag(tg)
	fmt.Printf(st[0])

	color.New(color.FgGreen).Printf("%s", st[1])
	fmt.Printf("%s", st[2])

}

func (p *parser) separateTag(tag []byte) [3]string {
	from := 1
	var newLine bool
	if tag[1] == '/' {
		from = 2
		newLine = true
	} else {
		p.NewLine = false
	}

	res := [3]string{}
	res[0] = fmt.Sprintf("%s", tag[0:from])

	to := from

	for ; to < len(tag)-1; to++ {
		if tag[to] == ' ' {
			break
		}
	}

	res[1] = fmt.Sprintf("%s", tag[from:to])
	p.Tags = append(p.Tags, res[1])
	if newLine {
		res[2] = fmt.Sprintf("%s\n", tag[to:])
	} else {
		res[2] = fmt.Sprintf("%s", tag[to:])
	}

	return res
}
