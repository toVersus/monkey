package evaluator

import (
	"testing"

	"github.com/toversus/monkey/lexer"
	"github.com/toversus/monkey/object"
	"github.com/toversus/monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.want)
	}
}

// testEval turns input into AST and evaluates them.
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, want int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%#+v)", obj, obj)
		return false
	}
	if result.Value != want {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, want)
		return false
	}

	return true
}
