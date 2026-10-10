package evaluator

import (
	"cardist/parser"
	"cardist/scanner"
)

// RuntimeError carries the operator token for line number reporting
type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func (e *RuntimeError) Error() string {
	return e.Message
}

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// Main evaluation entry point
func (e *Evaluator) Evaluate(expr parser.Expr) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			if runtimeErr, ok := r.(*RuntimeError); ok {
				err = runtimeErr
			} else {
				panic(r) // re-throw unknown panics
			}
		}
	}()

	return expr.Accept(e), nil
}

// 1. Literals
func (e *Evaluator) VisitLiteral(expr *parser.LiteralExpr) any {
	return expr.Value
}

// 2. Groupings
func (e *Evaluator) VisitGrouping(expr *parser.GroupingExpr) any {
	return expr.Expression.Accept(e)
}

// 3. Unary Operators (! and -)
func (e *Evaluator) VisitUnary(expr *parser.UnaryExpr) any {
	right := expr.Right.Accept(e)

	switch expr.Operator.Type {
	case scanner.TOKEN_MINUS:
		e.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case scanner.TOKEN_BANG:
		return !e.isTruthy(right)
	}

	return nil
}

// 4. Binary Operators (+, -, *, /, comparisons, equality)
func (e *Evaluator) VisitBinary(expr *parser.BinaryExpr) any {
	// Post-order traversal: left evaluated first, then right
	left := expr.Left.Accept(e)
	right := expr.Right.Accept(e)

	switch expr.Operator.Type {
	// Arithmetic
	case scanner.TOKEN_MINUS:
		e.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case scanner.TOKEN_STAR:
		e.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case scanner.TOKEN_SLASH:
		e.checkNumberOperands(expr.Operator, left, right)
		if right.(float64) == 0 {
			panic(&RuntimeError{Token: expr.Operator, Message: "Division by zero."})
		}
		return left.(float64) / right.(float64)
	case scanner.TOKEN_PLUS:
		// Overloaded: numbers or string concatenation
		if l, ok1 := left.(float64); ok1 {
			if r, ok2 := right.(float64); ok2 {
				return l + r
			}
		}
		if l, ok1 := left.(string); ok1 {
			if r, ok2 := right.(string); ok2 {
				return l + r
			}
		}
		panic(&RuntimeError{Token: expr.Operator, Message: "Operands must be two numbers or two strings."})

	// Comparisons
	case scanner.TOKEN_GREATER:
		e.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case scanner.TOKEN_GREATER_EQUAL:
		e.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case scanner.TOKEN_LESSER:
		e.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case scanner.TOKEN_LESSER_EQUAL:
		e.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)

	// Equality
	case scanner.TOKEN_EQUAL_EQUAL:
		return e.isEqual(left, right)
	case scanner.TOKEN_BANG_EQUAL:
		return !e.isEqual(left, right)
	}

	return nil
}

// --- Helper Functions ---

func (e *Evaluator) isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return true // Everything other than false and nil is truthy
}

func (e *Evaluator) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a == b
}

func (e *Evaluator) checkNumberOperand(operator scanner.Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(&RuntimeError{Token: operator, Message: "Operand must be a number."})
}

func (e *Evaluator) checkNumberOperands(operator scanner.Token, left, right any) {
	_, ok1 := left.(float64)
	_, ok2 := right.(float64)
	if ok1 && ok2 {
		return
	}
	panic(&RuntimeError{Token: operator, Message: "Operands must be numbers."})
}
