[![Go 1.23.1](https://img.shields.io/badge/go-1.23.1-9cf.svg)](https://golang.org/dl/)

# monkey-lang

TODO: write README documenting language & interpreter/compiler implementation

Monkey Programming Language, Interpreter, Compiler, & Virtual Machine written in Go. Inspired by [Writing An Interpreter in Go](https://interpreterbook.com/) & [Writing a Compiler in Go](https://compilerbook.com/) by [Thorsten Ball](https://thorstenball.com/).

## The Monkey Programming Language

Monkey is a programming language designed to help teach programming language theory & design, interpreters, and compilers.

![Monkey Logo](./docs/assets/monkey-lang.png)

Read more about Monkey at the [official site](https://monkeylang.org/).

## Binary

You can generate and interact with the Monkey binary on the REPL via `go build -o monkey && ./monkey`.

## Benchmarking

The `benchmark/` directory implements a Fibonacci benchmark of the interpreter & compiler/VM engines. To generate the benchmarking binary:

`go build -o fibonacci-benchmark ./benchmark`

To run the benchmark binary on each engine:

`./fibonacci-benchmark -engine=eval`

`./fibonacci-benchmark -engine=vm`
