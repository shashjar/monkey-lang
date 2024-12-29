package vm

import (
	"fmt"
	"math"
	"monkey/ast"
	"monkey/object"
)

func isNumerical(objectType object.ObjectType) bool {
	return objectType == object.INTEGER_OBJ || objectType == object.FLOAT_OBJ
}

func getNumericalValue(obj object.Object) (float64, bool) {
	switch obj := obj.(type) {
	case *object.Integer:
		return float64(obj.Value), false
	case *object.Float:
		return obj.Value, true
	default:
		panic(fmt.Sprintf("unsupported numerical type: %s", obj.Type()))
	}
}

func floatEquality(float1 float64, float2 float64) bool {
	return math.Abs(float1-float2) <= ast.FLOAT_64_EQUALITY_THRESHOLD
}
