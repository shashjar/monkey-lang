package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, statement := range program.Statements {
		if isMacroDefinition(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}

	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		definitionIndex := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

func ExpandMacros(program *ast.Program, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node { return macroExpansionModifier(node, env) })
}

func isMacroDefinition(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStatement.Value.(*ast.MacroLiteral)
	return ok
}

func addMacro(node ast.Statement, env *object.Environment) {
	letStatement, _ := node.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Body:       macroLiteral.Body,
		Env:        env,
	}

	env.Set(letStatement.Name.Value, macro)
}

func macroExpansionModifier(node ast.Node, env *object.Environment) ast.Node {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return node
	}

	macro, ok := isMacroCall(callExpression, env)
	if !ok {
		return node
	}

	args := quoteArgs(callExpression)
	evalEnv := extendMacroEnv(macro, args)

	evaluated := Eval(macro.Body, evalEnv)

	quote, ok := evaluated.(*object.Quote)
	if !ok {
		panic("we only support returning AST-nodes from macros")
	}

	return quote.Node
}

func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}
	for _, a := range exp.Arguments {
		args = append(args, &object.Quote{Node: a})
	}
	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	if len(args) != len(macro.Parameters) {
		panic(fmt.Sprintf("wrong number of arguments provided to macro. expected=%d, received=%d", len(macro.Parameters), len(args)))
	}

	env := object.NewEnclosedEnvironment(macro.Env)
	for i, param := range macro.Parameters {
		env.Set(param.Value, args[i])
	}
	return env
}
