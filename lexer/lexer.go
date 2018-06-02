package lexer

import "github.com/toversus/monkey/token"

// Lexer is used to take source code as input and output the tokens that represent the source code.
// TODO: Fully support Unicode (and emojis).
type Lexer struct {
	input string

	// position is current position in input (points to current char).
	position int

	// readPosition is current reading position in input (next to the current char).
	// It is useful to "peek" what comes up next after current char.
	readPosition int

	// ch is current char under examination, corresponding to the char in the position.
	// Change of the type to rune will be required if handling the full Unicode and UTF-8 chars.
	ch byte
}

// readChar throws the next char and advances the position in the input string.
// It only supports ASCII characters and doesn't aim to support full Unicode range
// due to remaining the simplicity.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// Set the ASCII code for the "NUL" char.
		// This means either "could not reead any chars yet" or "end of file".
		l.ch = 0
	} else {
		// Set the next char, but only valid in ASCII world.
		l.ch = l.input[l.readPosition]
	}
	// Advances the position in the input string.
	l.position = l.readPosition
	l.readPosition++
}

// New initialized the Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken returns a token parsed after examination of the current char
// and advances the pointers to the next char in input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// TODO: Consider to replace the branching method from switch to map.
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()
	return tok
}

// nextToken initializes the token passing through.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
