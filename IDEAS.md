# Ideas for Extensions/Improvements to Monkey

This is a living document storing some ideas for extensions & improvements that I might make to the Monkey programming language, interpreter, evaluator, compiler, and/or virtual machine implementations.

## Language Features

- [ ] Support floating-point numbers (currently only integers are supported)
- [ ] Add support for `else if` within conditional expressions
- [ ] Support `const` binding declarations in addition to `let`
- [ ] Boolean operations (`&&`, `||`)
- [ ] Modulo operator (`%`)
- [ ] Additional comparison operators (`>=` & `<=`)
- [ ] String operations: split, comparison with `==` & `!=`
- [ ] The lexer (`lexer/lexer.go`) currently only supports ASCII characters. Maybe extend this to Unicode (see p. 19-20 in WAIIG).
- [ ] Add support for macros into the compiler/VM engine (supported in interpreter but not yet compiler/VM)
- [ ] Array operations: join, `+` operator
- [ ] Add support for `switch` statements

# Builtin Functions

- [x] Add support for hashmaps in the `len` builtin function
- [ ] Add support for more language built-ins (`evaluator/builtins.go`), e.g. list comprehensions from Python - map, sum, filter, etc.

## REPL Features

- [x] When opening a REPL, ability to select interpreter/evaluator or compiler/VM engine
- [ ] Implement better printing of functions to the console (currently looks like: `Closure[0x140000ce160]`)
- [ ] Ability to navigate back and forth in REPL input with arrow keys & do multi-line input (look into [readline package](https://github.com/chzyer/readline))

## Additional Features

- [ ] Ability to execute Monkey programs from files (`.mo` extensions?) - again, be able to choose engine
- [ ] Ability to include comments in Monkey code
- [ ] Better error messages that point to line/column numbers for problematic tokens, both at compile-time and run-time
- [ ] Maybe write a UI to interact with Monkey - `wadackel` has written a great example of [this](https://github.com/wadackel/rs-monkey-lang)

## Compiler Internals

- [ ] Rely on an intermediate representation (IR) between the AST and bytecode, potentially to improve performance or simplify operations

## Benchmarking

- [ ] Can potentially implement a different/additional type of benchmark for comparing the interpreter/evaluator to compiler/VM. The `benchmark` directory currently only uses Fibonacci to benchmark the two engines against each other.

## Housekeeping

- [ ] Keep [README](README.md) documenting language features & implementation up-to-date
- [ ] In the AST modification functionality (`ast/modify.go`), implement thorough error-checking
