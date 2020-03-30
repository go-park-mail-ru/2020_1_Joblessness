package searchInterfaces

import "errors"

var (
	ErrUnknownRequest       = errors.New("invalid search parameters")
)