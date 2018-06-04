package parser

import (
	"github.com/toversus/monkey/ast"
	"github.com/toversus/monkey/lexer"
	"github.com/toversus/monkey/token"
)

// Parser is used to construct AST.
type Parser struct {
	l *lexer.Lexer

	// curToken is "pointers" (position and readPosition) to the current token.
	curToken token.Token
	// peekToken is also "pointers" to the next token.
	peekToken token.Token
}

// New initiates parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken is a helper function to advance both curToken and peekToken.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram constructs nodes of AST by recursive descent parser.
func (p *Parser) ParseProgram() *ast.Program {
	// Construct root node of the AST.
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// parseLetStatement constructs node with the token.LET token
// and advances the tokens until valid calls to expect peek.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: skip the expressions until encountering a semicolon.
	if !p.expectPeek(token.ASSIGN) {
		p.nextToken()
	}

	return stmt
}

// curTokenIs is useful method when fleshing out the parser.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs is useful method when fleshing out the parser.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek is assertion function to enforce the correctness of the order of tokens
// by checking the type of the next token.
// It returns nil and gets ignored if encountering a token that's not of the expected type,
// but this behavior makes it tough for debugging.
// TODO: Introduce error handling to the parser.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}
