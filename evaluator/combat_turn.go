package evaluator

import "cardist/parser"

func (e *Evaluator) executeTurn(stmt *parser.TurnStmt) error {
	if e.turnEnergy != nil {
		return &RuntimeError{Token: stmt.Keyword, Message: "Can't start a turn inside another turn."}
	}
	energy := float64(3)
	e.turnEnergy = &energy
	defer func() {
		e.turnEnergy = nil
		// Globals can be rebound in scripts, so reset the original combatants.
		for _, entity := range e.combatants {
			entity.block = 0
		}
	}()
	return e.executeBlock(stmt.Body, NewEnvironment(e.env))
}

func (e *Evaluator) activeEnergy() *float64 {
	if e.turnEnergy == nil {
		panic(&RuntimeError{Token: e.callSite, Message: "Energy is only available during a turn."})
	}
	return e.turnEnergy
}

func init() {
	natives["energy"] = &nativeFn{arity: 0, fn: func(e *Evaluator, _ []any) any {
		return *e.activeEnergy()
	}}

	natives["spend"] = &nativeFn{arity: 1, fn: func(e *Evaluator, args []any) any {
		energy := e.activeEnergy()
		amount := e.combatWholeNumber(args[0], "Amount")
		if amount > *energy {
			return false
		}
		*energy -= amount
		return true
	}}
}
