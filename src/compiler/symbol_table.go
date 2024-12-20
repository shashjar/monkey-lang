package compiler

// Represents a scope for some symbol/binding in a Monkey program.
type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
	BuiltInScope  SymbolScope = "BUILTIN"
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

	FreeSymbols []Symbol
}

func NewSymbolTable() *SymbolTable {
	store := make(map[string]Symbol)
	freeSymbols := []Symbol{}
	return &SymbolTable{
		store:       store,
		FreeSymbols: freeSymbols,
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

func (st *SymbolTable) DefineFunctionName(name string) Symbol {
	sym := Symbol{Name: name, Scope: FunctionScope, Index: 0}
	st.store[name] = sym
	return sym
}

func (st *SymbolTable) DefineBuiltIn(index int, name string) Symbol {
	sym := Symbol{Name: name, Scope: BuiltInScope, Index: index}
	st.store[name] = sym
	return sym
}

func (st *SymbolTable) defineFreeVar(original Symbol) Symbol {
	st.FreeSymbols = append(st.FreeSymbols, original)
	symbol := Symbol{Name: original.Name, Scope: FreeScope, Index: len(st.FreeSymbols) - 1}
	st.store[original.Name] = symbol
	return symbol
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.store[name]
	if !ok && st.outer != nil {
		sym, ok = st.outer.Resolve(name)
		if !ok {
			return sym, ok
		}

		if sym.Scope == GlobalScope || sym.Scope == BuiltInScope {
			return sym, ok
		}

		freeVar := st.defineFreeVar(sym)
		return freeVar, true
	}
	return sym, ok
}
