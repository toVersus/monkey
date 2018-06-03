package token

// Whispace is not here because it just acts as a separator for other tokens.
// Unlike Python, no interest towards the length of whitespace for the lexer (tokenizer).
const (
	// ILLEGAL signifies unknown token or character.
	ILLEGAL = "ILLEGAL"
	// EOF is used to pass on "end of file" to the parser.
	EOF = "EOF"

	// IDENT represents identifiers such as add, foobar, x, y, ...
	IDENT = "IDENT"
	// INT represents integer such as 123456.
	INT = "INT"

	// ASSIGN is used when binding some values to a name.
	ASSIGN = "="
	// PLUS represents to add left and right side of operator.
	// It is commonly known as infix notation.
	PLUS = "+"
	// MINUS represents to subtract left and right side of operator.
	// It is commonly known as infix notation.
	MINUS = "-"
	// BANG represents logically "not" operator.
	BANG = "!"
	// ASTERISK represents multiplication.
	ASTERISK = "*"
	// SLASH represents division remainder.
	SLASH = "/"

	// LT represents "less than".
	LT = "<"
	// GT represents "greater than".
	GT = ">"

	// COMMA is delimiter to separate multiple values.
	COMMA = ","
	// SEMICOLON is delimiter to represent the end of statement.
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// FUNCTION declares the definition of function.
	// In Monkey, functions are just values like integers or strings (first class functions),
	// and also can take other functions as arguments (higher order functions).
	FUNCTION = "FUNCTION"
	// LET binds the value (left side of '=' operator) to the name (right side of '=' operator).
	LET = "LET"
)

// TokenType is used to distinguish between different type of tokens.
// It is set to 'string' becase it allows us to debug easily without boilertemplate or helper function.
// This desicion degrates the performance of lexer compared to the case to use 'int' or 'byte' instead.
type TokenType string

// Token represents the data structure of Token.
// Type attribute is used to distinguish between 'integers' and 'right bracket'
// Literal attribute memorizes whether a 'number' token is a 5 or a 10,
// and this information will be reused in AST flow.
type Token struct {
	Type    TokenType
	Literal string
}

// keywords is the table of reserved keywords in language and its tokentype.
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent checks whether the given identifier is a reserved keyword or user-defined identifier.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
