package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/toversus/monkey/object"

	"github.com/toversus/monkey/compiler"
	"github.com/toversus/monkey/lexer"
	"github.com/toversus/monkey/parser"
	"github.com/toversus/monkey/vm"
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
// and passes it to an instance of lexer and compiles and executes the program
// until the end of source code.
func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)

	symbolTables := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTables.DefineBuiltin(i, v.Name)
	}

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

		comp := compiler.NewWithState(symbolTables, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		// update the constants reference because the compiler uses append internally
		// and previously allocated constants slice is included in the bytecode.
		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
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
