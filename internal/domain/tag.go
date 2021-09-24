package domain

import (
	"fmt"

	"github.com/tty2/xq/internal/domain/symbol"
)

type Tag struct {
	Bytes      []byte
	Name       string
	Attributes map[string]string
}

type attribute struct {
	Name  []byte
	Value []byte
}

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

	t.Name = fmt.Sprint(t.Bytes[startName:endName])

	return nil
}

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

	t.Name = fmt.Sprint(t.Bytes[startName:endName])

	var insideTag bool
	t.Attributes = map[string]string{}
	attr := attribute{}
	for i := endName; i < len(t.Bytes)-1; i++ {
		if insideTag {
			if symbol.IsQuote(t.Bytes[i]) {
				if len(attr.Name) != 0 {
					t.Attributes[fmt.Sprint(attr.Name)] = fmt.Sprint(attr.Value)

					attr = attribute{}
					insideTag = false
				}

				continue
			}
			attr.Value = append(attr.Value, t.Bytes[i])

			continue
		}
		if t.Bytes[i] != '=' {
			continue
		}
		if symbol.IsQuote(t.Bytes[i]) && t.Bytes[i-1] != '=' {
			insideTag = true

			continue
		}
		if t.Bytes[i] != ' ' {
			attr.Name = append(attr.Name, t.Bytes[i])
		}
	}

	return nil
}

func (t *Tag) GetAttributeValue(name string) (string, error) {
	v, ok := t.Attributes[name]
	if !ok {
		return "", fmt.Errorf("there is no attribute with name %s", name)
	}

	return v, nil
}
