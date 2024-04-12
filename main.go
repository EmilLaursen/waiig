package main

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/EmilLaursen/wiig/eval"
	"github.com/EmilLaursen/wiig/object"
	"github.com/EmilLaursen/wiig/parser"
	"github.com/EmilLaursen/wiig/repl"
)

func startRepl() {
	fmt.Printf("HELLO")
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

func runFiles(stdout, stderr io.Writer, files ...string) {
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("read file: %s, %s\n", file, err)
			return
		}

		p := parser.FromInput(string(data))
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			fmt.Fprintf(stderr, "Errors in file: %s\n", file)
			for _, msg := range p.Errors() {
				io.WriteString(stderr, msg)
				io.WriteString(stderr, "\n")
			}
			return
		}
		env := object.NewEnv()
		val := eval.Eval(program, env)
		io.WriteString(stdout, val.Inspect())
		io.WriteString(stdout, "\n")
	}
}

const usage string = `Usage:
%s repl

%s run [ FILES... ]
`

func main() {
	n := len(os.Args)
	fmt.Printf("args: %+v\n", os.Args)
	if n < 1 {
		fmt.Printf(usage, os.Args[0], os.Args[0])
		os.Exit(1)
	}
	switch os.Args[1] {
	case "repl":
		startRepl()
	case "run":
		runFiles(os.Stdout, os.Stderr, os.Args[2:]...)
	default:
		fmt.Printf(usage, os.Args[0], os.Args[0])
		os.Exit(1)
	}
}
