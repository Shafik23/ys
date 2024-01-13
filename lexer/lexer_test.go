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

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2, 3];
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
