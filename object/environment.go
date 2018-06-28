package object

// Environment is used to keep track of value by associating them with a name.
// It looks up in the outer scope if something is not found in the inner scope.
// The outer scope encloses the inner scope, otherwise the inner scope extends the outer one.
type Environment struct {
	store map[string]Object

	// outer represents enclosing environment.
	outer *Environment
}

// NewEncloseEnvironment makes enclosed environment.
func NewEncloseEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// Get also checks the enclosing environment for the given name.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
