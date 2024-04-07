package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/EmilLaursen/wiig/eval"
	"github.com/EmilLaursen/wiig/object"
	"github.com/EmilLaursen/wiig/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		p := parser.FromInput(line)
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		val := eval.Eval(program, env)
		if val != nil {
			io.WriteString(out, val.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t")
		io.WriteString(out, msg)
		io.WriteString(out, "\n")
	}
}
