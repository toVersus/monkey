package evaluator

import (
	"testing"

	"github.com/toversus/monkey/object"
)

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			`quote(unquote(4))`,
			`4`,
		},
		{
			`quote(unquote(4 + 4))`,
			`8`,
		},
		{
			`quote(8 + unquote(4 + 4))`,
			`(8 + 8)`,
		},
		{
			`quote(unquote(4 + 4) + 8)`,
			`(8 + 8)`,
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote. got=%T (%+v)",
				evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatal("quote.Node is nil")
		}

		if quote.Node.String() != test.want {
			t.Errorf("not equal. got=%q, want=%q",
				quote.Node.String(), test.want)
		}
	}
}
