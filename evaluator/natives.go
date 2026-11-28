package evaluator

import "time"

// nativeFn is a Callable that runs Go code instead of walking an AST body.
// Because it implements the same Callable as UserFunction, VisitCall checks
// arity and "not callable" once for both kinds.
type nativeFn struct {
	arity int
	fn    func(e *Evaluator, args []any) any
}

func (n *nativeFn) Arity() int { return n.arity }

func (n *nativeFn) Call(e *Evaluator, args []any) any { return n.fn(e, args) }

// natives is the table NewEvaluator copies into the global environment.
// Adding a native later = one more entry (combat.go adds its own in init).
var natives = map[string]*nativeFn{
	// clock() -> number: wall-clock seconds. Only compare two readings in tests.
	"clock": {arity: 0, fn: func(*Evaluator, []any) any {
		return float64(time.Now().UnixNano()) / 1e9
	}},
	// str(x) -> string: the same text print would show, via FormatValue.
	"str": {arity: 1, fn: func(_ *Evaluator, args []any) any {
		return FormatValue(args[0])
	}},
}
