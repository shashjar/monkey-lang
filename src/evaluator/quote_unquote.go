package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func quote(args []ast.Expression, env *object.Environment) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments provided to 'quote'. expected=1, got=%d", len(args))
	}

	node := evalUnquoteCalls(args[0], env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(node ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(node, func(node ast.Node) ast.Node { return unquoteEvalModifier(node, env) })
}

func unquoteEvalModifier(node ast.Node, env *object.Environment) ast.Node {
	if !isUnquoteCall(node) {
		return node
	}

	unquoteCall, _ := node.(*ast.CallExpression)
	if len(unquoteCall.Arguments) != 1 {
		return node
	}

	unquoteEval := Eval(unquoteCall.Arguments[0], env)
	return convertObjectToASTNode(unquoteEval)
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return callExpression.Function.TokenLiteral() == "unquote"
}
