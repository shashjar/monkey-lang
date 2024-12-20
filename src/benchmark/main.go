package main

import (
	"flag"
	"fmt"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"time"
)

var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

var input = `
	let fibonacci = fn(x) {
		if (x == 0) {
			return 0;
		} else {
			if (x == 1) {
				return 1;
			} else {
				fibonacci(x - 1) + fibonacci(x - 2);
			}
		}
	};

	fibonacci(35);
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if *engine == "vm" {
		// Compilation
		compiler := compiler.NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			fmt.Printf("compiler error: %s", err)
			return
		}

		// Execution on Virtual Machine (with timing/benchmarking)
		vm := vm.NewVM(compiler.Bytecode())
		startTime := time.Now()
		err = vm.Run()
		if err != nil {
			fmt.Printf("vm error: %s", err)
		}
		duration = time.Since(startTime)
		result = vm.LastPoppedStackElem()
	} else {
		env := object.NewEnvironment()
		startTime := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(startTime)
	}

	fmt.Printf(
		"engine=%s, result=%s, duration=%s\n",
		*engine,
		result.Inspect(),
		duration,
	)
}
