package repl

import (
	"fmt"
	"io"
	"monkey/compiler"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/vm"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

const PROMPT = ">> "
const CONTINUE_PROMPT = "... "

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

// The REPL for the Monkey programming language.
type REPL struct {
	out         io.Writer
	rl          *readline.Instance
	constants   []object.Object
	symbolTable *compiler.SymbolTable
	globals     []object.Object
}

// Creates a new REPL for the user to interact with Monkey at the top level.
func NewREPL(out io.Writer) (*REPL, error) {
	// Configure readline with custom settings
	config := &readline.Config{
		Prompt:            PROMPT,
		HistoryFile:       "/tmp/monkey_repl_history.tmp",
		HistoryLimit:      200,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return nil, fmt.Errorf("error initializing REPL: %s", err)
	}

	constants := []object.Object{}
	symbolTable := compiler.NewSymbolTable()
	for i, builtin := range object.BuiltIns {
		symbolTable.DefineBuiltIn(i, builtin.Name)
	}
	globals := make([]object.Object, vm.GlobalsSize)

	return &REPL{
		out:         out,
		rl:          rl,
		constants:   constants,
		symbolTable: symbolTable,
		globals:     globals,
	}, nil
}

// Starts the REPL for the Monkey programming language compiler & VM for the user to interact with.
func (r *REPL) Start() {
	defer r.rl.Close()

	for {
		// Read input
		input, err := r.readMultiLineInput()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			return
		} else if err != nil {
			fmt.Fprintf(r.out, "Error reading input on REPL: %s\n", err)
			continue
		}

		// Handle exit command
		if strings.TrimSpace(input) == "exit" {
			return
		}

		// Skip empty inputs
		if strings.TrimSpace(input) == "" {
			continue
		}

		r.executeInput(input)
	}
}

func (r *REPL) ExecuteFile(filename string) {
	defer r.rl.Close()

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(r.out, "Error reading from file: %s\n", err)
		return
	}

	input := string(bytes)
	r.executeInput(input)
}

func (r *REPL) readMultiLineInput() (string, error) {
	var lines []string
	r.rl.SetPrompt(PROMPT)

	for {
		line, err := r.rl.Readline()
		if err != nil {
			return "", err
		}

		// Detect multi-line input (i.e. backslash at the end of the line)
		if !strings.HasSuffix(line, "\\") {
			lines = append(lines, line)
			break
		}

		// Remove backslash for proper input handling
		line = strings.TrimSuffix(line, "\\")
		lines = append(lines, line)
		r.rl.SetPrompt(CONTINUE_PROMPT)
	}

	return strings.Join(lines, "\n"), nil
}

func (r *REPL) executeInput(input string) {
	// Lexing
	l := lexer.NewLexer(input)

	// Parsing
	p := parser.NewParser(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(r.out, p.Errors())
		return
	}

	// Compilation
	compiler := compiler.NewCompilerWithState(r.symbolTable, r.constants)
	err := compiler.Compile(program)
	if err != nil {
		fmt.Fprintf(r.out, "Whoops! Compilation failed:\n %s\n", err)
	}

	bytecode := compiler.Bytecode()
	r.constants = bytecode.Constants

	// Virtual Machine (VM)
	vm := vm.NewVMWithGlobalsStore(bytecode, r.globals)
	err = vm.Run()
	if err != nil {
		fmt.Fprintf(r.out, "Whoops! Executing bytecode failed:\n %s\n", err)
		return
	}

	// Printing Output
	lastPopped := vm.LastPoppedStackElem()
	if lastPopped != nil {
		io.WriteString(r.out, lastPopped.Inspect())
		io.WriteString(r.out, "\n")
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
