package lexer

import (
	"testing"
)

// TestLexer tests the tokenization of the input.
func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:     "Simple arithmetic",
			input:    "10 + 25 - 8",
			expected: []Token{INT, PLUS, INT, MINUS, INT},
		},
		{
			name:     "Single plus",
			input:    "+",
			expected: []Token{PLUS},
		},
		{
			name:     "Single minus",
			input:    "-",
			expected: []Token{MINUS},
		},
		{
			name:     "Whitespace only",
			input:    "   ",
			expected: []Token{EOF},
		},
		{
			name:     "Empty input",
			input:    "",
			expected: []Token{EOF},
		},
		{
			name:     "Multiple operators",
			input:    "+-+-+-",
			expected: []Token{PLUS, MINUS, PLUS, MINUS, PLUS, MINUS},
		},
		{
			name:     "Numbers with whitespace",
			input:    "10   20  30",
			expected: []Token{INT, INT, INT},
		},
		{
			name:     "Mixed whitespace",
			input:    "10 +\n25 -   8",
			expected: []Token{INT, PLUS, INT, MINUS, INT},
		},
		{
			name:     "Newlines",
			input:    "\n\n\n",
			expected: []Token{EOF},
		},
		{
			name:     "Complex expression",
			input:    "10+20-30+40-50",
			expected: []Token{INT, PLUS, INT, MINUS, INT, PLUS, INT, MINUS, INT},
		},
		{
			name:     "Identifiers",
			input:    "x y z",
			expected: []Token{IDENT, IDENT, IDENT},
		},
		{
			name:     "Keywords",
			input:    "fn let if else return",
			expected: []Token{FUNCTION, LET, IF, ELSE, RETURN},
		},
		{
			name:     "Operators",
			input:    "= == ! != < >",
			expected: []Token{ASSIGN, EQ, BANG, NOT_EQ, LT, GT},
		},
		{
			name:     "Delimiters",
			input:    ", ; ( ) { }",
			expected: []Token{COMMA, SEMICOLON, LPAREN, RPAREN, LBRACE, RBRACE},
		},
		{
			name:     "Mixed code",
			input:    "let x = 5; if (x == 5) { x = x * 10; }",
			expected: []Token{LET, IDENT, ASSIGN, INT, SEMICOLON, IF, LPAREN, IDENT, EQ, INT, RPAREN, LBRACE, IDENT, ASSIGN, IDENT, ASTERISK, INT, SEMICOLON, RBRACE},
		},
		{
			name:     "Function declaration",
			input:    "fn add(x, y) { return x + y; }",
			expected: []Token{FUNCTION, IDENT, LPAREN, IDENT, COMMA, IDENT, RPAREN, LBRACE, RETURN, IDENT, PLUS, IDENT, SEMICOLON, RBRACE},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)

			for j, expected := range tt.expected {
				tok, _ := l.NextToken()

				// Skip whitespace tokens.
				for tok == WS {
					tok, _ = l.NextToken()
				}

				if tok != expected {
					t.Fatalf("%s - tokens[%d] wrong. expected=%q, got=%q",
						tt.name, j, expected, tok)
				}
			}

			// Check for EOF at the end.
			if tok, _ := l.NextToken(); tok != EOF {
				t.Fatalf("%s - tokens at the end. expected EOF, got=%q", tt.name, tok)
			}
		})
	}
}
