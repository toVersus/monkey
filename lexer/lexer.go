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

	l.skipWhitespace()

	// TODO: Consider to replace the branching method from switch to map.
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
		// Token will be updated when detecting "==" tokens.
		// The following process will be abstracted when supporting more two-char tokens.
		if l.peekChar() == '=' {
			// memorize current char before readChar calls overwrites current char.
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = newToken(token.BANG, l.ch)
		// Token will be updated when detecting "!=" tokens.
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOTEQ, Literal: literal}
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}
		// If reaching end of this block, the current char cannot be handled.
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

// readIdentifier reads in an identifier and advances lexer's position
// until it encounters a non-letter char.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter is helper function to check whether the given argument is a letter or not.
// Changes to this function will heavily impact on the language itself.
// Currently, snake case representation is supported for the identifier and keyword,
// but both '!' and '?' are not recognized as identifier.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// nextToken initializes the token passing through.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// skipWhitespace skips over whitespace, tabspace and newline.
// This rule depends on the language being lexed.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readNumber reads in an identifier and advances lexer's positionuntil it encounters a non-numeric char.
// readNumber is almost same implementation as readIdentifier except for its usage of isDigit instead of isLetter.
// But, for simplicity and ease of understanding, these two function are not generalize
// by passing in the char-idenrifying function as arguments.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isDigit is helper function to check whether the given argument is integer or not.
// This means that it doesn't support float, numbers in hex and octal notation in this stage.
// TODO: Support float or other digit format.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// peekChar only looks ahead in the input and grasps what to be reuturned after readChar call in advance.
// It doesn't move around in it.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// readString reads character by calling readChar until encountering either a closing double quote or end of the input.
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
