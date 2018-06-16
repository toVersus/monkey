package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/toversus/monkey/lexer"
	"github.com/toversus/monkey/parser"
)

// PROMPT is used in the prompt of Monkey interactive mode.
const PROMPT = ">> "

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
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
