package evaluator

import (
	"cardist/parser"
	"fmt"
)

// maxCallDepth is far below where Go's stack runs out, and far above what
// ordinary recursion such as fib(20) needs.
const maxCallDepth = 1000

// Callable is anything a call expression can invoke. Week 3's natives
// implement it too, so arity and "not callable" are checked in one place.
type Callable interface {
	Arity() int
	Call(e *Evaluator, args []any) any
}

// UserFunction is a `func` declaration plus the scope it was declared in.
type UserFunction struct {
	decl    *parser.FunctionStmt
	closure *Environment
}

func (f *UserFunction) Arity() int { return len(f.decl.Params) }

// Call runs the body in a new scope enclosing the CLOSURE (where the function
// was written), not the caller's scope. That one choice makes Cardist lexical.
func (f *UserFunction) Call(e *Evaluator, args []any) any {
	env := NewEnvironment(f.closure)
	for i, param := range f.decl.Params {
		env.Define(param.Lexeme, args[i])
	}
	if err := e.executeBlock(f.decl.Body, env); err != nil {
		panic(err) // Evaluate's recover turns this back into a runtime error
	}
	return nil // shortcut: no `return` yet, so every call yields nil; Week 3 adds results
}

func (e *Evaluator) VisitCall(expr *parser.CallExpr) any {
	callee := expr.Callee.Accept(e)
	args := make([]any, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = arg.Accept(e)
	}

	fn, err := callee.(Callable)
	if !err {
		panic(&RuntimeError{Token: expr.Paren, Message: "Can only call functions."})
	}
	if len(args) != fn.Arity() {
		panic(&RuntimeError{Token: expr.Paren, Message: fmt.Sprintf(
			"Expected %d arguments but got %d.", fn.Arity(), len(args))})
	}
	if e.callDepth >= maxCallDepth {
		panic(&RuntimeError{Token: expr.Paren, Message: "Stack overflow."})
	}
	e.callDepth++
	defer func() { e.callDepth-- }() // also runs when an error unwinds through here
	return fn.Call(e, args)
}
