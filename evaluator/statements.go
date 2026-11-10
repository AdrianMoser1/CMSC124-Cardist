package evaluator

import (
	"cardist/parser"
	"fmt"
)

// ExecuteScript runs a file's statements in order; only print produces output.
func (e *Evaluator) ExecuteScript(script *parser.Script) error {
	for _, statement := range script.Statements {
		if err := e.execute(statement); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteREPL is ExecuteScript, except a top-level `expr;` prints its value.
func (e *Evaluator) ExecuteREPL(script *parser.Script) error {
	for _, statement := range script.Statements {
		if exprStmt, ok := statement.(*parser.ExpressionStmt); ok {
			statement = &parser.PrintStmt{Expression: exprStmt.Expression}
		}
		if err := e.execute(statement); err != nil {
			return err
		}
	}
	return nil
}

func (e *Evaluator) execute(statement parser.Stmt) error {
	switch stmt := statement.(type) {
	case *parser.PrintStmt:
		value, err := e.Evaluate(stmt.Expression)
		if err == nil && e.Emit != nil {
			e.Emit(value)
		}
		return err

	case *parser.ExpressionStmt:
		_, err := e.Evaluate(stmt.Expression)
		return err

	case *parser.VarStmt:
		var value any // `var x;` starts as nil
		if stmt.Initializer != nil {
			var err error
			if value, err = e.Evaluate(stmt.Initializer); err != nil {
				return err
			}
		}
		e.env.Define(stmt.Name.Lexeme, value)
		return nil

	case *parser.WhileStmt:
		for {
			condition, err := e.Evaluate(stmt.Condition)
			if err != nil {
				return err
			}
			if !e.isTruthy(condition) {
				return nil
			}
			if err := e.execute(stmt.Body); err != nil {
				return err
			}
		}

	case *parser.IfStmt:
		condition, err := e.Evaluate(stmt.Condition)
		if err != nil {
			return err
		}
		if e.isTruthy(condition) {
			return e.execute(stmt.ThenBranch)
		}
		if stmt.ElseBranch != nil {
			return e.execute(stmt.ElseBranch)
		}
		return nil

	case *parser.BlockStmt:
		// Enclosing = where the block is written, which makes scope lexical.
		return e.executeBlock(stmt.Statements, NewEnvironment(e.env))
	}
	return fmt.Errorf("Runtime error: unsupported statement type %T", statement)
}

func (e *Evaluator) executeBlock(statements []parser.Stmt, blockEnv *Environment) error {
	previous := e.env
	e.env = blockEnv
	// Restore even on error, or the REPL keeps running inside a dead block scope.
	defer func() { e.env = previous }()

	for _, statement := range statements {
		if err := e.execute(statement); err != nil {
			return err
		}
	}
	return nil
}
