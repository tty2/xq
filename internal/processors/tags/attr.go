package tags

import "github.com/tty2/xq/internal/domain/symbol"

func pickAttributesNames(tag []byte) []string {
	if len(tag) == 0 || tag[0] != symbol.OpenBracket {
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
