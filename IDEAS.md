# Ideas for Extensions/Improvements to Monkey

TODO: this is a living document storing some ideas for extensions & improvements that I might make to the Monkey programming language, interpreter, compiler, and/or virtual machine implementations.

## Language Features

- The lexer (`lexer/lexer.go`) currently only supports ASCII characters. Maybe extend this to Unicode (see p. 19-20 in WAIIG).
- Add support for more language built-ins (`evaluator/builtins.go`), e.g. list comprehensions from Python - map, sum, filter, etc.
- Add support for `else if` within conditional expressions
- Support `const` binding declarations in addition to `let`

## Compiler Internals

- Rely on an intermediate representation (IR) between the AST and bytecode, potentially to improve performance or simplify operations

## Benchmarking

- Can potentially implement a different/additional type of benchmark for comparing the interpreter to compiler/VM. The `benchmark` directory currently only uses Fibonacci to benchmark the two engines against each other.

## Housekeeping

- In the AST modification functionality (`ast/modify.go`), implement thorough error-checking
