package vm

import (
	"monkey/object"
	"testing"
)

func TestPuts(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `puts("")`,
			expected: Null,
		},
		{
			input:    `puts("hi there!")`,
			expected: Null,
		},
		{
			input:    `puts("hello", "world!")`,
			expected: Null,
		},
		{
			input:    `puts(["writing", "tests", "is", "fun"])`,
			expected: Null,
		},
		{
			input:    `puts({"test": "this"})`,
			expected: Null,
		},
		{
			input:    `puts()`,
			expected: &object.Error{Message: "need at least one argument provided to `puts`"},
		},
	}

	runVMTests(t, tests)
}

func TestLen(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `len("")`,
			expected: 0,
		},
		{
			input:    `len("three")`,
			expected: 5,
		},
		{
			input:    `len("hello world")`,
			expected: 11,
		},
		{
			input:    `len([])`,
			expected: 0,
		},
		{
			input:    `len([1, 2, 3])`,
			expected: 3,
		},
		{
			input:    `len({})`,
			expected: 0,
		},
		{
			input:    `len({1: 2, "hi": "there", true: false})`,
			expected: 3,
		},
		{
			input:    `len()`,
			expected: &object.Error{Message: "wrong number of arguments. expected=1, got=0"},
		},
		{
			input:    `len("one", "two")`,
			expected: &object.Error{Message: "wrong number of arguments. expected=1, got=2"},
		},
		{
			input:    `len(42)`,
			expected: &object.Error{Message: "argument to `len` is not supported, got INTEGER"},
		},
		{
			input:    `len(true)`,
			expected: &object.Error{Message: "argument to `len` is not supported, got BOOLEAN"},
		},
	}

	runVMTests(t, tests)
}

func TestFirst(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `first([1, 2, 3])`,
			expected: 1,
		},
		{
			input:    `first([true, {1: 2, "a": "b"}, 10, false, -6])`,
			expected: true,
		},
		{
			input:    `first([])`,
			expected: &object.Error{Message: "array is empty; no first element"},
		},
		{
			input:    `first()`,
			expected: &object.Error{Message: "wrong number of arguments. expected=1, got=0"},
		},
		{
			input:    `first([1, 2], [3, 4])`,
			expected: &object.Error{Message: "wrong number of arguments. expected=1, got=2"},
		},
		{
			input:    `first("hi")`,
			expected: &object.Error{Message: "argument to `first` is not supported, got STRING"},
		},
	}

	runVMTests(t, tests)
}

func TestLast(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `last([40, 41, 42])`,
			expected: 42,
		},
		{
			input:    `last([])`,
			expected: &object.Error{Message: "array is empty; no last element"},
		},
		{
			input:    `last()`,
			expected: &object.Error{Message: "wrong number of arguments. expected=1, got=0"},
		},
		{
			input:    `last([1, 2], [3, 4])`,
			expected: &object.Error{Message: "wrong number of arguments. expected=1, got=2"},
		},
		{
			input:    `last(6)`,
			expected: &object.Error{Message: "argument to `last` is not supported, got INTEGER"},
		},
	}

	runVMTests(t, tests)
}

func TestRest(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `rest([])`,
			expected: Null,
		},
		{
			input:    `rest([1])`,
			expected: []int{},
		},
		{
			input:    `rest([1, 2, 3])`,
			expected: []int{2, 3},
		},
		{
			input:    `rest(6)`,
			expected: &object.Error{Message: "argument to `rest` is not supported, got INTEGER"},
		},
	}

	runVMTests(t, tests)
}

func TestAppend(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `append([], 1)`,
			expected: []int{1},
		},
		{
			input:    `append([4, -7], 1)`,
			expected: []int{4, -7, 1},
		},
		{
			input:    `append([4, 5])`,
			expected: &object.Error{Message: "wrong number of arguments. expected=2, got=1"},
		},
	}

	runVMTests(t, tests)
}

func TestJoin(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `join([])`,
			expected: "",
		},
		{
			input:    `join(["hello"])`,
			expected: "hello",
		},
		{
			input:    `join(["hello", " world", "!"])`,
			expected: "hello world!",
		},
		{
			input:    `join(["hello", "world", "i", "am", "here"], " ")`,
			expected: "hello world i am here",
		},
		{
			input:    `join(["testing", "different", "delimiter"], ":")`,
			expected: "testing:different:delimiter",
		},
		{
			input:    `join([4, 5], [5, 6], [6, 7])`,
			expected: &object.Error{Message: "wrong number of arguments. expected 1 or 2, got=3"},
		},
		{
			input:    `join(["hi", "there"], 6)`,
			expected: &object.Error{Message: "expected delimiter passed to `join` to be a string, got INTEGER"},
		},
		{
			input:    `join([1, 2, 3], " ")`,
			expected: &object.Error{Message: "elements of array passed to `join` must be strings, got INTEGER"},
		},
		{
			input:    `join(" ", ",")`,
			expected: &object.Error{Message: "first argument to `join` must be an array, got STRING"},
		},
	}

	runVMTests(t, tests)
}
