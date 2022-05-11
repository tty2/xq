/*
Package symbol keeps all the symbols are used by parser.
*/
package symbol

// Symbols are used by parser.
const (
	CloseBracket = '>'
	OpenBracket  = '<'

	NewLine        = 10 // '\n'
	CarriageReturn = 13 // '\r'

	Quote       = 39 // '
	DoubleQuote = 34 // "

	Space = 32
)

// IsQuote checks if symbols is a single or double quote.
func IsQuote(s byte) bool {
	return s == Quote || s == DoubleQuote
}

func IsSpace(s byte) bool {
	if s == ' ' || s == '\n' || s == '\t' {
		return true
	}

	return false
}
