package processor

import (
	"errors"
	"fmt"

	"github.com/tty2/xq/internal/domain/symbol"
)

func pickAttributesNames(tag []byte) []string {
	if len(tag) < 3 || tag[0] != symbol.OpenBracket || tag[len(tag)-1] != symbol.CloseBracket {
		return nil
	}

	var i int
	// skip tag name
	for ; i < len(tag) && tag[i] != ' '; i++ {
	}

	var isAttrName bool
	attrs := []string{}
	attrName := []byte{}
	quotes := []byte{}
	for ; i < len(tag); i++ {
		if symbol.IsQuote(tag[i]) {
			if len(quotes) == 0 {
				quotes = append(quotes, tag[i])

				continue
			}

			if len(quotes) > 0 && tag[i] == quotes[0] { // close quote
				quotes = []byte{}
			} else {
				quotes = append(quotes, tag[i])
			}
		}

		if len(quotes) > 0 {
			continue
		}

		if tag[i] == '=' {
			attrs = append(attrs, string(attrName))
			isAttrName = false
			attrName = []byte{}

			continue
		}
		if isAttrName {
			attrName = append(attrName, tag[i])

			continue
		}
		if tag[i] == ' ' { // space is followed by attribute name
			isAttrName = true
		}
	}

	return attrs
}

func pickAttributeValue(targetName string, tag []byte) (string, error) {
	if len(tag) < 3 || tag[0] != symbol.OpenBracket || tag[len(tag)-1] != symbol.CloseBracket {
		return "", errors.New("invalid tag")
	}

	var i int
	// skip tag name
	for ; i < len(tag) && tag[i] != ' '; i++ {
	}

	var isAttrName bool
	var writeValue bool
	attrName := []byte{}
	targetValue := []byte{}
	quotes := []byte{}
	for ; i < len(tag); i++ {
		if symbol.IsQuote(tag[i]) {
			if len(quotes) == 0 {
				quotes = append(quotes, tag[i])

				continue
			}

			if len(quotes) > 0 && tag[i] == quotes[0] { // close quote
				quotes = []byte{}
			} else {
				quotes = append(quotes, tag[i])
			}
		}

		if len(quotes) > 0 {
			if writeValue {
				targetValue = append(targetValue, tag[i])
			}

			continue
		}

		if writeValue {
			return string(targetValue), nil
		}

		if tag[i] == '=' {
			if string(attrName) == targetName {
				writeValue = true
			}
			isAttrName = false
			attrName = []byte{}

			continue
		}
		if isAttrName {
			attrName = append(attrName, tag[i])

			continue
		}
		if tag[i] == ' ' { // space is followed by attribute name
			isAttrName = true
		}
	}

	return "", fmt.Errorf("there is no attribute with name `%s`", targetName)
}
