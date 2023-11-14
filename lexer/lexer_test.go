package lexer

import (
	"testing"

	"github.com/shafik23/ys/token"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "Simple arithmetic",
			input: "10 + 25 - 8",
			expected: []token.Token{
				{Type: token.INT, Literal: "10"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.INT, Literal: "25"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.INT, Literal: "8"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "Variable assignment",
			input: "x = 10",
			expected: []token.Token{
				{Type: token.IDENT, Literal: "x"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "10"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "Boolean expressions",
			input: "true == false",
			expected: []token.Token{
				{Type: token.TRUE, Literal: "true"},
				{Type: token.EQ, Literal: "=="},
				{Type: token.FALSE, Literal: "false"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "Grouped expressions",
			input: "(5 + 5) * 2",
			expected: []token.Token{
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "5"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.INT, Literal: "5"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.ASTERISK, Literal: "*"},
				{Type: token.INT, Literal: "2"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "Function definition",
			input: "fn add(a, b) { return a + b; }",
			expected: []token.Token{
				{Type: token.FUNCTION, Literal: "fn"},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "a"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "b"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.IDENT, Literal: "a"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.IDENT, Literal: "b"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "Control structures",
			input: "if (5 < 10) { return true; } else { return false; }",
			expected: []token.Token{
				{Type: token.IF, Literal: "if"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.INT, Literal: "5"},
				{Type: token.LT, Literal: "<"},
				{Type: token.INT, Literal: "10"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.TRUE, Literal: "true"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.ELSE, Literal: "else"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.FALSE, Literal: "false"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "String literals",
			input: `"Hello, World!"`,
			expected: []token.Token{
				{Type: token.STRING, Literal: "Hello, World!"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "Comments",
			input: "// This is a comment\nx = 10",
			expected: []token.Token{
				{Type: token.COMMENT, Literal: "// This is a comment"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "10"},
				{Type: token.EOF, Literal: ""},
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)

			for i, expected := range tt.expected {
				tok := l.NextToken()

				if tok.Type != expected.Type {
					t.Fatalf("%s - tokens[%d] wrong type. expected=%q, got=%q",
						tt.name, i, expected.Type, tok.Type)
				}

				if tok.Literal != expected.Literal {
					t.Fatalf("%s - tokens[%d] wrong literal. expected=%q, got=%q",
						tt.name, i, expected.Literal, tok.Literal)
				}
			}

			tok := l.NextToken()
			if tok.Type != token.EOF {
				t.Fatalf("%s - token at the end. expected EOF, got=%+v", tt.name, tok)
			}
		})
	}
}
