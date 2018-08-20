package compiler

// SymbolScope represents unique scopes to differentiate between different scopes.
type SymbolScope string

const (
	LocalScope  SymbolScope = "LOCAL"
	GlobalScope SymbolScope = "GLOBAL"
)

// Symbol all the necessary information to be used to retrieve
// by a given identifier.
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable is a data structure used to interpreters and compilers
// to associate identifiers with information.
type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

// Define adds the symbol defined in global scope to the symbol table and returns symbol.
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	s.store[name] = symbol

	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve retrieves the symbol table by name and returns symbol.
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		return obj, ok
	}
	return obj, ok
}
