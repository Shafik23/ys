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
			expected: []Token{INTEGER, PLUS, INTEGER, MINUS, INTEGER},
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
			expected: []Token{INTEGER, INTEGER, INTEGER},
		},
		{
			name:     "Mixed whitespace",
			input:    "10 +\n25 -   8",
			expected: []Token{INTEGER, PLUS, INTEGER, MINUS, INTEGER},
		},
		{
			name:     "Newlines",
			input:    "\n\n\n",
			expected: []Token{EOF},
		},
		{
			name:     "Complex expression",
			input:    "10+20-30+40-50",
			expected: []Token{INTEGER, PLUS, INTEGER, MINUS, INTEGER, PLUS, INTEGER, MINUS, INTEGER},
		},
		// Add more complex test cases here as you expand your lexer's capabilities.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)

			for j, expected := range tt.expected {
				tok, _ := l.Next()

				// Skip whitespace tokens.
				for tok == WS {
					tok, _ = l.Next()
				}

				if tok != expected {
					t.Fatalf("%s - tokens[%d] wrong. expected=%q, got=%q",
						tt.name, j, expected, tok)
				}
			}

			// Check for EOF at the end.
			if tok, _ := l.Next(); tok != EOF {
				t.Fatalf("%s - tokens at the end. expected EOF, got=%q", tt.name, tok)
			}
		})
	}
}
