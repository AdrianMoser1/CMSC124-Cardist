package evaluator

import "cardist/scanner"

// Environment holds one scope's bindings. enclosing links to the surrounding
// scope and is nil for the global environment.
type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]any),
		enclosing: enclosing,
	}
}

// Define always writes to this scope and never walks outward, which is what
// lets an inner declaration shadow an outer one. Redeclaring a name in the
// same scope overwrites it.
func (env *Environment) Define(name string, value any) {
	env.values[name] = value
}

// Get returns the value from the innermost scope that holds the name.
func (env *Environment) Get(name scanner.Token) any {
	if value, ok := env.values[name.Lexeme]; ok {
		return value
	}
	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}
	panic(&RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."})
}

// Assign updates an existing binding. It never creates one, so a misspelled
// name is an error instead of a new variable.
func (env *Environment) Assign(name scanner.Token, value any) {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return
	}
	if env.enclosing != nil {
		env.enclosing.Assign(name, value)
		return
	}
	panic(&RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."})
}
