package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
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

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}
