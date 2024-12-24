package ast

// Represents a function that optionally modifies an AST node in some way.
type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {

	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)

	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)

	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)

	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)

	case *BlockStatement:
		for i, statement := range node.Statements {
			node.Statements[i], _ = Modify(statement, modifier).(Statement)
		}

	case *IfExpression:
		modifiedClauses := []ConditionalClause{}
		for _, clause := range node.Clauses {
			modifiedClauseCondition, _ := Modify(clause.Condition, modifier).(Expression)
			modifiedClauseConsequence, _ := Modify(clause.Consequence, modifier).(*BlockStatement)
			modifiedClauses = append(modifiedClauses, ConditionalClause{Condition: modifiedClauseCondition, Consequence: modifiedClauseConsequence})
		}
		node.Clauses = modifiedClauses
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}

	case *FunctionLiteral:
		for i := range node.Parameters {
			node.Parameters[i], _ = Modify(node.Parameters[i], modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)

	case *ArrayLiteral:
		for i, element := range node.Elements {
			node.Elements[i], _ = Modify(element, modifier).(Expression)
		}

	case *HashMapLiteral:
		newKVPairs := make(map[Expression]Expression)
		for key, val := range node.KVPairs {
			newKey, _ := Modify(key, modifier).(Expression)
			newVal, _ := Modify(val, modifier).(Expression)
			newKVPairs[newKey] = newVal
		}
		node.KVPairs = newKVPairs

	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)

	}

	return modifier(node)
}
