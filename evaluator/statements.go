package evaluator

import (
	"cardist/parser"
	"fmt"
)

// ExecuteScript runs statements in source order.
// Only PrintStmt sends a value to Emit.
func (e *Evaluator) ExecuteScript(script *parser.Script) error {
	for _, statement := range script.Statements {
		var expression parser.Expr
		shouldPrint := false

		switch stmt := statement.(type) {
		case *parser.PrintStmt:
			expression = stmt.Expression
			shouldPrint = true

		case *parser.ExpressionStmt:
			expression = stmt.Expression

		case *parser.VarStmt:
			// `var x;` has no initializer, so x starts as nil.
			var value any
			if stmt.Initializer != nil {
				initial, err := e.Evaluate(stmt.Initializer)
				if err != nil {
					return err
				}
				value = initial
			}
			e.env.Define(stmt.Name.Lexeme, value)
			continue

		default:
			return fmt.Errorf(
				"Runtime error: unsupported statement type %T",
				statement,
			)
		}

		value, err := e.Evaluate(expression)
		if err != nil {
			return err
		}

		if shouldPrint && e.Emit != nil {
			e.Emit(value)
		}
	}

	return nil
}
