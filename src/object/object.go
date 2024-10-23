package object

import "fmt"

// Represents a type of internally-represented object in the Monkey programming language.
type ObjectType string

const (
	NULL_OBJ    = "NULL"
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
)

// Represents an object in the Monkey programming language.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Represents null.
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

// Represents an integer.
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Represents a boolean.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
