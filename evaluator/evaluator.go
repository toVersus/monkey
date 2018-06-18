package evaluator

import (
	"github.com/toversus/monkey/ast"
	"github.com/toversus/monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
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

	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
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

// nativeBoolToBooleanObject converts native bool object to reference of "true" and "false" instances
// instead of allocating new object.
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
