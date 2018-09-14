package evaluator

import (
	"github.com/toversus/monkey/ast"
	"github.com/toversus/monkey/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
