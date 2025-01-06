package object

import "fmt"

func IsNumerical(objectType ObjectType) bool {
	return objectType == INTEGER_OBJ || objectType == FLOAT_OBJ
}

func GetNumericalValue(obj Object) (float64, bool, error) {
	switch obj := obj.(type) {
	case *Integer:
		return float64(obj.Value), false, nil
	case *Float:
		return obj.Value, true, nil
	default:
		return 0, false, fmt.Errorf("unsupported numerical type: %s", obj.Type())
	}
}
