package domain

import (
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
	Name  []byte
	Value []byte
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
