package domain

import (
	"github.com/tty2/xq/internal/domain/color"
	"github.com/tty2/xq/internal/domain/symbol"
)

// Tag represents a tag object.
// It make sense to work with is with `Bytes` field populated only otherwise all the methods
// will not work.
// `Bytes` field MUST start from open bracket `<` and finish with close bracket `>` and
// MUST be 3 symbols or bigger otherwise tag is invalid.
type Tag struct {
	Bytes      []byte
	Name       string
	Attributes map[string]string
}

type attribute struct {
	Name        []byte
	Value       []byte
	Quote       byte
	NextIsQuote bool
	InsideValue bool
}

// Validate checks if tag `Bytes` has 3 or more symbols, starts with open bracket
// and the last symbol is close bracket.
func (t *Tag) Validate() error {
	if len(t.Bytes) < 3 {
		return ErrTagShort
	}

	if t.Bytes[0] != symbol.OpenBracket {
		return ErrTagInvalidStart
	}

	if t.Bytes[len(t.Bytes)-1] != symbol.CloseBracket {
		return ErrTagInvalidEnd
	}

	return nil
}

// SetName takes name from tag and set to `Name` field.
func (t *Tag) SetName() error {
	err := t.Validate()
	if err != nil {
		return err
	}

	startName := 1         // name starts after open bracket
	if t.Bytes[1] == '/' { // closed tag
		startName = 2
	}

	endName := startName

	for ; endName < len(t.Bytes)-1; endName++ {
		if t.Bytes[endName] == ' ' {
			break
		}
	}

	t.Name = string(t.Bytes[startName:endName])

	return nil
}

// SetNameAndAttributes takes name from tag and set it to `Name` field +
// takes all attributes and their values and put them to `Attributes`.
func (t *Tag) SetNameAndAttributes() error {
	err := t.Validate()
	if err != nil {
		return err
	}

	startName := 1         // name starts after open bracket
	if t.Bytes[1] == '/' { // closed tag
		startName = 2
	}

	endName := startName

	for ; endName < len(t.Bytes)-1; endName++ {
		if t.Bytes[endName] == ' ' {
			break
		}
	}

	t.Name = string(t.Bytes[startName:endName])

	if startName == 2 {
		return nil // it's forbidden to set attributes to close tag
	}

	var insideTag bool
	t.Attributes = map[string]string{}
	attr := attribute{}
	for i := endName; i < len(t.Bytes)-1; i++ {
		if insideTag {
			if symbol.IsQuote(t.Bytes[i]) {
				if len(attr.Name) != 0 {
					t.Attributes[string(attr.Name)] = string(attr.Value)

					attr = attribute{}
					insideTag = false
				}

				continue
			}
			attr.Value = append(attr.Value, t.Bytes[i])

			continue
		}
		if t.Bytes[i] == '=' {
			continue
		}
		if symbol.IsQuote(t.Bytes[i]) && t.Bytes[i-1] == '=' {
			insideTag = true

			continue
		}
		if t.Bytes[i] != ' ' {
			attr.Name = append(attr.Name, t.Bytes[i])
		}
	}

	return nil
}

// ColorizeTag colorizes tag.
func ColorizeTag(tg []byte) []byte {
	ln := len(tg)

	coloredTag := make([]byte, 0, ln)

	startName := 1    // name starts after open bracket
	if tg[1] == '/' { // closed tag
		startName = 2
	}

	endName := startName
	for ; endName < len(tg)-1; endName++ {
		if tg[endName] == ' ' {
			break
		}
	}

	coloredTag = append(coloredTag, tg[:startName]...)        // add open bracket
	coloredTag = append(coloredTag, []byte(color.Red)...)     // add red color
	coloredTag = append(coloredTag, tg[startName:endName]...) // tag name

	attr := attribute{
		Value: []byte{},
	}
	for i := endName; i < ln-1; i++ {
		if attr.NextIsQuote {
			if symbol.IsQuote(tg[i]) {
				attr.Quote = tg[i]
				attr.NextIsQuote = false
				attr.InsideValue = true
			}

			continue
		}
		if attr.InsideValue {
			if tg[i] == attr.Quote && (len(attr.Value) == 0 || attr.Value[len(attr.Value)-1] != '\\') {
				attr.InsideValue = false
				coloredTag = append(coloredTag, symbol.Space)
				coloredTag = append(coloredTag, []byte(color.Green)...)
				coloredTag = append(coloredTag, attr.Name...)
				coloredTag = append(coloredTag, []byte(color.White)...)
				coloredTag = append(coloredTag, '=', attr.Quote)
				coloredTag = append(coloredTag, attr.Value...)
				coloredTag = append(coloredTag, attr.Quote)
				attr = attribute{
					Value: []byte{},
				}
			} else {
				attr.Value = append(attr.Value, tg[i])
			}

			continue
		}
		if tg[i] == ' ' {
			continue
		}
		if tg[i] == '=' { // value of attribute
			attr.NextIsQuote = true
			coloredTag = append(coloredTag, []byte(color.White)...)

			continue
			// i != ln-3 in order do not colorize `/` sign inside an empty tag in case like this `<...attr="value" />`
		} else if tg[i] == '/' && i == ln-2 { // end attribute value
			coloredTag = append(coloredTag, tg[i])

			continue
		}
		attr.Name = append(attr.Name, tg[i])
	}
	coloredTag = append(coloredTag, []byte(color.White)...)
	coloredTag = append(coloredTag, tg[ln-1]) // add close bracket

	return coloredTag
}
