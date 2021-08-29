/*
This part of processing is responsible for full document parsing with colorizing and
making indentation. This processing doesn't parse and separate tags, attributes and
values from document structure but optimized for fast processing with less memory
consumption.
*/
package main

import (
	"bytes"
	"fmt"
	"log"
)

func (p *parser) parseFullDocument(chunk []byte) {
	for i := range chunk {
		// skip carriage return and new line from data in order do not duplicate with created ones by parser
		if p.SkipData && (chunk[i] == ' ' || chunk[i] == '\t') {
			continue
		}
		if chunk[i] == newLine || chunk[i] == carriageReturn {
			p.SkipData = true

			continue
		}
		p.SkipData = false

		if p.InsideTag {
			p.CurrentTag.Bytes = append(p.CurrentTag.Bytes, chunk[i])

			if chunk[i] == closeBracket {
				p.CurrentTag.Brackets--

				if p.CurrentTag.Brackets > 0 {
					continue
				}

				p.InsideTag = false
				p.printTag()
				p.SkipData = true // skip if there are empty symbols beeween close tag and new data
			} else if chunk[i] == openBracket {
				p.CurrentTag.Brackets++
			}

			continue
		}

		if chunk[i] == openBracket {
			p.InsideTag = true
			p.CurrentTag = tag{
				Bytes: []byte{chunk[i]},
			}
			p.CurrentTag.Brackets++

			if len(p.Data) > 0 {
				// nolint forbidigo: printf is executed on purpose here
				fmt.Printf("%s\n", append(bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation), p.Data...))
				p.Data = []byte{}
			}

			continue
		}

		p.Data = append(p.Data, chunk[i])
	}
}

// nolint forbidigo: printf in this method is executed on purpose
func (p *parser) printTag() {
	if len(p.CurrentTag.Bytes) < minTagSize {
		log.Fatalf("tag size is too small = %d, tag is `%s`", len(p.CurrentTag.Bytes), p.CurrentTag.Bytes)
	}
	if p.CurrentTag.Bytes[1] == '!' || p.CurrentTag.Bytes[1] == '?' { // service tag, comment or cdata
		fmt.Printf("%s\n", append(bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation), p.CurrentTag.Bytes...))

		return
	}

	fmt.Printf("%s\n", p.colorizeTag())

	if p.CurrentTag.Bytes[len(p.CurrentTag.Bytes)-2] != '/' {
		p.Indentation++
	}
}

func (p *parser) downIndent() {
	p.Indentation--
}

func (p *parser) colorizeTag() []byte {
	ln := len(p.CurrentTag.Bytes)

	coloredTag := make([]byte, 0, p.Indentation+ln)

	startName := 1                    // name starts after open bracket
	if p.CurrentTag.Bytes[1] == '/' { // closed tag
		startName = 2
		p.Indentation--
		defer p.downIndent()
	}

	endName := startName
	for ; endName < len(p.CurrentTag.Bytes)-1; endName++ {
		if p.CurrentTag.Bytes[endName] == ' ' {
			break
		}
	}

	coloredTag = append(coloredTag, bytes.Repeat([]byte(" "), p.IndentItemSize*p.Indentation)...) // add indentation
	coloredTag = append(coloredTag, p.CurrentTag.Bytes[:startName]...)                            // add open bracket
	coloredTag = append(coloredTag, []byte(red)...)                                               // add red color
	coloredTag = append(coloredTag, p.CurrentTag.Bytes[startName:endName]...)                     // tag name

	attr := attribute{
		Value: []byte{},
	}
	for i := endName; i < ln-1; i++ {
		if attr.NextIsQuote {
			if isQuote(p.CurrentTag.Bytes[i]) {
				attr.Quote = p.CurrentTag.Bytes[i]
				attr.NextIsQuote = false
				attr.InsideValue = true
			}

			continue
		}
		if attr.InsideValue {
			if p.CurrentTag.Bytes[i] == attr.Quote && attr.Value[len(attr.Value)-1] != '\\' {
				attr.InsideValue = false
				coloredTag = append(coloredTag, space)
				coloredTag = append(coloredTag, []byte(green)...)
				coloredTag = append(coloredTag, attr.Name...)
				coloredTag = append(coloredTag, []byte(white)...)
				coloredTag = append(coloredTag, '=', attr.Quote)
				coloredTag = append(coloredTag, attr.Value...)
				coloredTag = append(coloredTag, attr.Quote)
				attr = attribute{
					Value: []byte{},
				}
			} else {
				attr.Value = append(attr.Value, p.CurrentTag.Bytes[i])
			}

			continue
		}
		if p.CurrentTag.Bytes[i] == ' ' {
			continue
		}
		if p.CurrentTag.Bytes[i] == '=' { // value of attribute
			attr.NextIsQuote = true
			coloredTag = append(coloredTag, []byte(white)...)

			continue
			// i != ln-3 in order do not colorize `/` sign inside an empty tag in case like this `<...attr="value" />`
		} else if p.CurrentTag.Bytes[i] == '/' && i == ln-2 { // end attribute value
			coloredTag = append(coloredTag, p.CurrentTag.Bytes[i])

			continue
		}
		attr.Name = append(attr.Name, p.CurrentTag.Bytes[i])
	}
	coloredTag = append(coloredTag, []byte(white)...)
	coloredTag = append(coloredTag, p.CurrentTag.Bytes[ln-1]) // add close bracket

	return coloredTag
}

func isQuote(s byte) bool {
	return s == quote || s == doubleQuote
}
