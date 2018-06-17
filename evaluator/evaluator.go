package evaluator

import (
	"github.com/toversus/monkey/ast"
	"github.com/toversus/monkey/object"
)

// Eval traverses the AST and evaluates basic types.
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements because Eval always starts from the top of the tree.
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions starts here
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}
