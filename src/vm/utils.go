package vm

import (
	"math"
	"monkey/ast"
)

func floatEquality(float1 float64, float2 float64) bool {
	return math.Abs(float1-float2) <= ast.FLOAT_64_EQUALITY_THRESHOLD
}
