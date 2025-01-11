# Ideas for Extensions/Improvements to Monkey

This is a living document storing some ideas for extensions & improvements that I might make to the Monkey programming language, interpreter, evaluator, compiler, and/or virtual machine implementations.

## Language Features

- [x] Support floating-point numbers (currently only integers are supported)
- [x] Add support for `else if` within conditional expressions
- [x] Allow for reassignment of values to variables declared with let
- [x] Support `const` binding declarations in addition to `let`
- [x] Logical boolean operators (`&&`, `||`)
- [x] Modulus operator (`%`)
- [x] Add `//` (integer division) as a separate operator from `/`
- [x] Support `+=`, `-=`, `*=`, `/=`, `//=` operations
- [x] Add exponentiation operation with `**`
- [x] Additional comparison operators (`>=` & `<=`)
- [x] Strings: comparison with `==` & `!=`
- [x] Arrays: `+` operator
- [x] Implement `while` loops
- [x] Add basic support for `for` loops
- [ ] Add support for `break` and `continue` in `while` and `for` loops
- [ ] For `for` loops, the `Init`, `Condition`, and `Afterthought` expressions are all currently required - maybe allow for these to be optional
- [x] Add basic support for `switch` statements
- [ ] `switch` statements currently just use equality (`==`) for comparison - maybe allow for switching based on the type of some variable, like in Go
- [x] Maybe support postfix operators `++` and `--`
- [ ] The lexer (`lexer/lexer.go`) currently only supports ASCII characters. Maybe extend this to Unicode (see p. 19-20 in WAIIG).
- [ ] Add support for macros into the compiler/VM engine (supported in interpreter but not yet compiler/VM)

# Builtin Functions

- [x] Add support for hashmaps in the `len` builtin function
- [x] Strings: `split` operation
- [x] Arrays: `join` operation
- [x] Arrays: `sum`
- [ ] Arrays: `map`
- [ ] Arrays: `filter`

## REPL Features

- [x] When opening a REPL, ability to select interpreter/evaluator or compiler/VM engine
- [x] Implement better printing of functions to the console (currently looks like: `Closure[0x140000ce160]`)
- [x] Ability to navigate back and forth in REPL input with arrow keys & do multi-line input (look into [readline package](https://github.com/chzyer/readline))

## Additional Features

- [x] Ability to execute Monkey programs from files (`.mo` extensions?) - again, be able to choose engine
- [x] Better error messages that point to line/column numbers for problematic tokens, both at compile-time and run-time
- [ ] Ability to include comments in Monkey code
- [ ] Maybe write a UI to interact with Monkey - `wadackel` has written a great example of this: https://github.com/wadackel/rs-monkey-lang

## Compiler Internals

- [ ] Rely on an intermediate representation (IR) between the AST and bytecode, potentially to improve performance or simplify operations

## Benchmarking

- [ ] Can potentially implement additional benchmarks for comparing the interpreter/evaluator to compiler/VM: arithmetic operations, string concatenation, Sieve of Eratosthenes, etc.

## Housekeeping / Tech Debt

- [ ] Keep [README](README.md) documenting language features & implementation up-to-date
- [x] Maybe refactor some code files/split into more files (some files are getting pretty long)
- [x] In the `*ast.InfixExpression` handling in `compiler/compiler.go`, is it worth just adding an `OpLessThan` bytecode instruction so that this case doesn't have to be handled separately from the rest of the logic?
- [ ] In the AST modification functionality (`ast/modify.go`), implement thorough error-checking
