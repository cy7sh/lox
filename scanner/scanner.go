package scanner

import (
	"github.com/singurty/lox/token"
)

// Scanner transforms the source into tokens
type Scanner struct {
	source string
}

func New(source string) Scanner {
	scanner := Scanner{source: source}
	return scanner
}
