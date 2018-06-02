package lexer

import (
	"testing"

	"github.com/toversus/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		wantType    token.TokenType
		wantLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	// New will be implemented to return *Lexer.
	l := New(input)

	for i, test := range tests {
		tok := l.NextToken()

		if tok.Type != test.wantType {
			t.Fatalf("tests[%d] - wrong tokentype. want=%q, got=%q",
				i, test.wantType, tok.Type)
		}

		if tok.Literal != test.wantLiteral {
			t.Fatalf("tests[%d] - wrong literal. want=%q, got=%q",
				i, test.wantLiteral, tok.Literal)
		}
	}
}
