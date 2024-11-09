package compiler

// Represents a scope for some symbol/binding in a Monkey program.
type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

// Represents a symbol stored in the symbol table, associated with some scope.
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// Represents a symbol table associating Monkey identifiers with information.
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}

func (st *SymbolTable) Define(name string) Symbol {
	sym := Symbol{
		Name:  name,
		Scope: GlobalScope,
		Index: st.numDefinitions,
	}
	st.store[name] = sym
	st.numDefinitions += 1
	return sym
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.store[name]
	return sym, ok
}
