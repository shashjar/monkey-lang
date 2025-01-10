package main

import (
	"flag"
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

// By default, give the user the compiler/VM engine, but allow for interpreter/evaluator access if specified.
var engine = flag.String("engine", "vm", "use 'vm' or 'eval'")

// Entrypoint for the Monkey interpreter program.
func main() {
	flag.Parse()

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands!\n")

	if *engine == "vm" {
		r, err := repl.NewREPL(os.Stdout)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		r.Start()
	} else if *engine == "eval" {
		repl.StartInterpreter(os.Stdin, os.Stdout)
	} else {
		fmt.Printf("Invalid engine to use for REPL: %q\n", *engine)
		os.Exit(1)
	}
}
