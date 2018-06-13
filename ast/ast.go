package ast

import (
	"bytes"

	"github.com/toversus/monkey/token"
)

type Node interface {
	// TokenLiteral returns literal value of the associated token.
	// It is used only for debugging and testing.
	TokenLiteral() string

	// String prints AST nodes for debugging ans compares them with other AST nodes,
	// and it is handy in tests.
	String() string
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

// String returns whole program back as a string for readable tests.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement ...
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

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

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

func (i *Identifier) String() string { return i.Value }

// ReturnStatement consists of solely of the keyword "return" and an "expression".
type ReturnStatement struct {
	Token token.Token

	// ReturnValue is the expression to be returned.
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement is wrapper for adding it to the Statements slice of ast.Program
// to reuse work on the parser.
type ExpressionStatement struct {
	// The first token of the expression.
	Token token.Token

	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token token.Token
	// Operator only contains two types of operator, "-" and "!".
	Operator string
	// Right represents the expression to the right of the operator.
	Right Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression is almost same structure as PrefixExpression,
// but it holds expression in left side of infix operator.
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }
