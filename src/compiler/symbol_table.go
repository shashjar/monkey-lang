package compiler

// Represents a scope for some symbol/binding in a Monkey program.
type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

// Represents a symbol stored in the symbol table, associated with some scope.
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// Represents a symbol table associating Monkey identifiers with information.
type SymbolTable struct {
	outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	st := NewSymbolTable()
	st.outer = outer
	return st
}

func (st *SymbolTable) Define(name string) Symbol {
	sym := Symbol{Name: name, Index: st.numDefinitions}
	if st.outer == nil {
		sym.Scope = GlobalScope
	} else {
		sym.Scope = LocalScope
	}

	st.store[name] = sym
	st.numDefinitions += 1
	return sym
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.store[name]
	if !ok && st.outer != nil {
		return st.outer.Resolve(name)
	}
	return sym, ok
}
