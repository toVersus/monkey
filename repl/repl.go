package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/toversus/monkey/evaluator"
	"github.com/toversus/monkey/lexer"
	"github.com/toversus/monkey/parser"
)

const (
	// PROMPT is used in the prompt of Monkey interactive mode.
	PROMPT = ">> "

	// MONKEYFACE is printed out along with the friendly error messages.
	MONKEYFACE = `             __,__
    .--.  .-"     "-.  .--.
   / .. \/  .-. .-.  \/ .. \
  | |  '|  /   Y   \  |'  | |
  | \   \  \ 0 | 0 /  /   / |
   \ '- ,\.-"""""""-./, -' /
    ''-' /_   ^ ^   _\ '-''
        |  \._   _./  |
        \   \ '~' /   /
         '._ '-=-' _.'
            '-----'
`
)

// Start reads from input source code until encountering a newline
// and passes it to an instance of lexer and prints all the tokens
// until the end of source code.
func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		if !sc.Scan() {
			return
		}

		line := sc.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEYFACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
