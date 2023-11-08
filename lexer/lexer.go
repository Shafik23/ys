// File: lexer/lexer.go

package lexer

import (
	"unicode"
	"unicode/utf8"
)

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Symbols
	ASSIGN
	PLUS
	MINUS
	ASTERISK
	SLASH
	BANG
	LT
	GT
	EQ
	NOT_EQ

	// Delimiters
	COMMA
	SEMICOLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE

	// Literals
	IDENT // identifiers
	INT   // integers
	STR   // string literals

	// Keywords
	FUNCTION
	LET
	IF
	ELSE
	RETURN
)

// define a map for keywords
var keywords = map[string]Token{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// Lexer represents a lexer.
type Lexer struct {
	input        string
	pos          int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
}

// NewLexer returns a new instance of Lexer.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Initialize the first char
	return l
}

// readChar gets the next character and advances our position in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" character, signifies we're at EOF
	} else {
		l.ch, _ = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}
	l.pos = l.readPosition
	l.readPosition++
}

// NextToken returns the next token and literal value.
func (l *Lexer) NextToken() (tok Token, literal string) {
	var lit string
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit = string(ch) + string(l.ch)
			tok = EQ
		} else {
			tok = ASSIGN
			lit = string(l.ch)
		}
	case ';':
		tok = SEMICOLON
		lit = string(l.ch)
	case '(':
		tok = LPAREN
		lit = string(l.ch)
	case ')':
		tok = RPAREN
		lit = string(l.ch)
	case ',':
		tok = COMMA
		lit = string(l.ch)
	case '+':
		tok = PLUS
		lit = string(l.ch)
	case '-':
		tok = MINUS
		lit = string(l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			lit = string(ch) + string(l.ch)
			tok = NOT_EQ
		} else {
			tok = BANG
			lit = string(l.ch)
		}
	case '/':
		tok = SLASH
		lit = string(l.ch)
	case '*':
		tok = ASTERISK
		lit = string(l.ch)
	case '<':
		tok = LT
		lit = string(l.ch)
	case '>':
		tok = GT
		lit = string(l.ch)
	case '{':
		tok = LBRACE
		lit = string(l.ch)
	case '}':
		tok = RBRACE
		lit = string(l.ch)
	case 0:
		tok = EOF
		lit = ""
	default:
		if isLetter(l.ch) {
			lit = l.readIdentifier()
			tok = lookupIdent(lit)
			return tok, lit
		} else if unicode.IsDigit(l.ch) {
			tok = INT
			lit = l.readNumber()
			return tok, lit
		} else {
			tok = ILLEGAL
			lit = string(l.ch)
		}
	}

	l.readChar()
	return tok, lit
}

// lookupIdent checks the keywords table to see whether the given identifier is in fact a keyword.
func lookupIdent(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// readIdentifier reads in an identifier and advances the lexer's positions until it encounters a non-letter-character.
func (l *Lexer) readIdentifier() string {
	position := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

// readNumber reads a number and advances the lexer's positions until it encounters a non-digit character.
func (l *Lexer) readNumber() string {
	position := l.pos
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

// skipWhitespace skips any whitespace characters in the input.
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// peekChar returns the next character without moving the position.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
		return ch
	}
}

// isLetter checks if the character is a letter or underscore.
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}
