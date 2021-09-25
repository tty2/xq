package domain

import "errors"

// Errors for export.
var (
	ErrTagShort        = errors.New("tag can't be less then 3 bytes")
	ErrTagInvalidStart = errors.New("tag must start from open bracket symbol")
	ErrTagInvalidEnd   = errors.New("tag must end with close bracket symbol")
)
