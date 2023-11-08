package lexer

import (
	"unicode"
)

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Symbols
	PLUS
	MINUS

	// Literals
	INTEGER
)

// The Lexer holds the state of the scanner.
type Lexer struct {
	input string // the string being scanned
	pos   int    // current position in the input
}

// NewLexer returns a new instance of Lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Next returns the next token and literal value.
func (l *Lexer) Next() (tok Token, literal string) {
	// Read the next rune.
	if l.pos >= len(l.input) {
		return EOF, ""
	}

	ch := l.input[l.pos]

	// Skip whitespace.
	if unicode.IsSpace(rune(ch)) {
		l.pos++
		return WS, " "
	}

	// Check if the rune is a known symbol.
	switch ch {
	case '+':
		l.pos++
		return PLUS, "+"
	case '-':
		l.pos++
		return MINUS, "-"
	}

	// If the rune is a digit, scan the entire number.
	if unicode.IsDigit(rune(ch)) {
		start := l.pos
		for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
			l.pos++
		}
		return INTEGER, l.input[start:l.pos]
	}

	// If we haven't returned by now, it's an illegal token.
	l.pos++
	return ILLEGAL, string(ch)
}
