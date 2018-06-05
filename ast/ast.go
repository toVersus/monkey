package ast

import "github.com/toversus/monkey/token"

type Node interface {
	// TokenLiteral returns literal value of the associated token.
	// It is used only for debugging and testing.
	TokenLiteral() string
}

// Statement doesn't produce a value.
type Statement interface {
	Node
	// statementNode is helper function for guiding the Go compiler
	// and causing it to throw errors when it is unused.
	statementNode()
}

// Expression produces a value.
type Expression interface {
	Node
	// expressionNode is helper function for guiding the Go compiler
	// and causing it to throw errors when it is unused.
	expressionNode()
}

// Program is the root node of every AST that parser produces.
// Every valid Monkey program is a series of statements.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

//
type LetStatement struct {
	// Token is the token.LET token.
	Token token.Token

	// Name is the identifier of the binding.
	Name *Identifier

	// Value for the expression produces the value.
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier is used to hold the identifier of the binding.
// It is not Expression when it binds a variable, e.g. let x = 5;
// But in other cases, it is Expression, e.g. let x = valueProducingIdentifier;
// All in one piece, Identifier is Expression because of keeping it simple,
// keeping the number of different node types small.
type Identifier struct {
	// Token is the token.IDENT token
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// ReturnStatement consists of solely of the keyword "return" and an "expression".
type ReturnStatement struct {
	Token token.Token

	// ReturnValue is the expression to be returned.
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
