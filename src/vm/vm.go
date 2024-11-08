package vm

import (
	"fmt"
	"monkey/bytecode"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

var Null = &object.Null{}

// Represents a virtual machine used to execute bytecode instructions generated by the Monkey programming language compiler.
type VM struct {
	constants    []object.Object
	instructions bytecode.Instructions

	stack []object.Object
	sp    int // Always points to the next free slot in the stack. The top of the stack is stack[sp - 1].
}

func NewVM(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := bytecode.Opcode(vm.instructions[ip])

		switch op {
		case bytecode.OpConstant:
			constIndex := bytecode.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case bytecode.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case bytecode.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case bytecode.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case bytecode.OpPop:
			vm.pop()
		case bytecode.OpJumpNotTruthy:
			jumpToPos := int(bytecode.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = jumpToPos - 1 // Set to `pos - 1` since this loop increments ip on each iteration
			}
		case bytecode.OpJump:
			jumpToPos := int(bytecode.ReadUint16(vm.instructions[ip+1:]))
			ip = jumpToPos - 1 // Set to `pos - 1` since this loop increments ip on each iteration

		case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case bytecode.OpEqual, bytecode.OpNotEqual, bytecode.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case bytecode.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}
		case bytecode.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("invalid opcode received: %d", op)
		}
	}

	return nil
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp += 1

	return nil
}

func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp -= 1
	return obj
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	} else {
		return False
	}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func (vm *VM) executeBinaryOperation(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeBinaryIntegerOperation(op bytecode.Opcode, left object.Object, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case bytecode.OpAdd:
		result = leftValue + rightValue
	case bytecode.OpSub:
		result = leftValue - rightValue
	case bytecode.OpMul:
		result = leftValue * rightValue
	case bytecode.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown binary integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeComparison(op bytecode.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	} else if leftType == object.BOOLEAN_OBJ && rightType == object.BOOLEAN_OBJ {
		return vm.executeBooleanComparison(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary comparison: %s %s", leftType, rightType)
}

func (vm *VM) executeIntegerComparison(op bytecode.Opcode, left object.Object, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case bytecode.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue == rightValue))
	case bytecode.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue != rightValue))
	case bytecode.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown binary integer comparison operator: %d", op)
	}
}

func (vm *VM) executeBooleanComparison(op bytecode.Opcode, left object.Object, right object.Object) error {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value
	switch op {
	case bytecode.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue == rightValue))
	case bytecode.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue != rightValue))
	default:
		return fmt.Errorf("unknown binary boolean comparison operator: %d", op)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}
