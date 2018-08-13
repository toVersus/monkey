package compiler

// SymbolScope represents unique scopes to differentiate between different scopes.
type SymbolScope string

const (
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
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

// Define adds the symbol defined in global scope to the symbol table and returns symbol.
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve retrieves the symbol table by name and returns symbol.
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return obj, ok
}
