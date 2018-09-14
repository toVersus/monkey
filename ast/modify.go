package ast

type ModifyFunc func(Node) Node

func Modify(node Node, modifier ModifyFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	}

	return modifier(node)
}
