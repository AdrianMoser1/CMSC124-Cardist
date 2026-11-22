package evaluator

import (
	"cardist/parser"
	"cardist/scanner"
	"fmt"
)

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
	err := e.executeBlock(f.decl.Body, env)
	if ret, isReturn := err.(*returnSignal); isReturn {
		return ret.value
	}
	if err != nil {
		panic(err) // Evaluate's recover turns this back into a runtime error
	}
	return nil // fell off the end of the body
}

// returnSignal carries a `return` value out of the function body. It travels
// as an error so loops and blocks stop and executeBlock still restores scope.
// If nothing catches it, the return was outside any function.
type returnSignal struct {
	keyword scanner.Token
	value   any
}

func (r *returnSignal) Error() string {
	return fmt.Sprintf("[line %d] Runtime error: Can't return from top-level code.", r.keyword.Line)
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
	return fn.Call(e, args)
}
