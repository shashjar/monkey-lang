package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
)

const PROMPT = ">> "

const MONKEY_FACE = `            __,__
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

// Starts the REPL for the Monkey programming language compiler & VM for the user to interact with.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		// Reading Input
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Lexing
		line := scanner.Text()
		l := lexer.NewLexer(line)

		// Parsing
		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// Compilation
		compiler := compiler.NewCompiler()
		err := compiler.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Whoops! Compilation failed:\n %s\n", err)
		}

		// Virtual Machine (VM)
		vm := vm.NewVM(compiler.Bytecode())
		err = vm.Run()
		if err != nil {
			fmt.Fprintf(out, "Whoops! Executing bytecode failed:\n %s\n", err)
		}

		// Printing Output
		stackTop := vm.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

// Starts the REPL for the Monkey programming language interpreter for the user to interact with.
func StartInterpreter(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		// Reading Input
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Lexing
		line := scanner.Text()
		l := lexer.NewLexer(line)

		// Parsing
		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// Macro Expansion
		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		// Evaluation
		evaluated := evaluator.Eval(expanded, env)

		// Printing Output
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
