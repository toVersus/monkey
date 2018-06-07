package parser

import (
	"fmt"

	"github.com/toversus/monkey/ast"
	"github.com/toversus/monkey/lexer"
	"github.com/toversus/monkey/token"
)

const (
	// To assign increment numbers starting from 1 as the constant values.
	_ int = iota

	// The order and the relation to each other are critical for representing precedence.
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        //myFunction(X)
)

// Parser is used to construct AST.
type Parser struct {
	l *lexer.Lexer

	errors []string

	// curToken is "pointers" (position and readPosition) to the current token.
	curToken token.Token
	// peekToken is also "pointers" to the next token.
	peekToken token.Token

	// These maps can be used to get the correct parser function for the current token type.
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression

	// argument represents left side of the infix operator to be parsed.
	infixParseFn func(ast.Expression) ast.Expression
)

// New initiates parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

// parseIdentifier just returns token and its value.
// It doesn't advance the tokens by calling nextToken.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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
	case token.RETURN:
		return p.parseReturnStatement()

	// Workaround for passing TestLetStatements at this time.
	case token.SEMICOLON:
		return nil

	default:
		return p.parseExpressionStatement()
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
	if !p.curTokenIs(token.SEMICOLON) {
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
// It returns error expressions and shows the expected type of token if encountering a mismatch of type of token.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// Errors returns error messages.
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError is helper function to detect mismatch of the type of peekToken.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// parseReturnStatement just constructs ast.ReturnStatement with the current token.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: skip the expressions until encountering a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// registerPrefix registers prefix parsing function as associated token type.
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers infix parsing function as associated token type.
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// Semicolon is optional in this context.
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}