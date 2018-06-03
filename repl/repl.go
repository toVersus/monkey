package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/toversus/monkey/lexer"
	"github.com/toversus/monkey/token"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
