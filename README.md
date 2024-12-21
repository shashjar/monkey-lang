[![Go 1.23.1](https://img.shields.io/badge/go-1.23.1-9cf.svg)](https://golang.org/dl/)

# monkey-lang

Monkey Programming Language, Interpreter, Compiler, & Virtual Machine written in Go. Inspired by [Writing An Interpreter in Go](https://interpreterbook.com/) & [Writing a Compiler in Go](https://compilerbook.com/) by [Thorsten Ball](https://thorstenball.com/).

I'm actively extending this language implementation with new features, for which I have many [ideas](IDEAS.md).

## The Monkey Programming Language

Monkey is a programming language designed to help teach programming language theory & design, interpreters, and compilers.

![Monkey Logo](./docs/assets/monkey-lang.png)

Read more about Monkey at the [official site](https://monkeylang.org/).

## Usage

To use the Monkey REPL, run:

```
go build -o monkey && ./monkey
```

<img src="./docs/assets/monkey-usage.png" alt="Monkey Usage" style="box-shadow: 5px 5px 15px rgba(0,0,0,0.3); border-radius: 10px;">


By default, this runs the Monkey compiler & virtual machine, with entrypoint `main.go` into `repl/repl.go`. The `StartInterpreter` function can be used to instead spin up the interpreter-driven REPL.

## Benchmarks

The `benchmark/` directory implements a Fibonacci benchmark of the interpreter/evaluator & compiler/VM engines. To generate the benchmarking binary:

```
go build -o fibonacci-benchmark ./benchmark
```

To run the benchmark binary on each engine:

```
./fibonacci-benchmark -engine=eval
```

```
./fibonacci-benchmark -engine=vm
```

## Implementation Details

### Interpreter

### Compiler & Virtual Machine

## Language Documentation

### Table of Contents

### Summary

### Integers & Arithmetic Operations

### Booleans

### Comparison Operators

### Conditionals

### Bindings

### Strings

### Arrays

### Hashmaps

### Functions

### Built-In Functions
