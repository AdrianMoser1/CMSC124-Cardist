package evaluator

import (
	"fmt"
	"math"
)

// Entity holds combat stats and owns its draw pile and hand.
type Entity struct {
	name       string
	hp, block  float64
	deck, hand []*Card
}

func (en *Entity) String() string {
	return fmt.Sprintf("<%s hp=%s block=%s>", en.name, formatNumber(en.hp), formatNumber(en.block))
}

func defineCombatants(globals *Environment) [2]*Entity {
	// shortcut: starting stats and deck are fixed, upgrade with definition blocks.
	player := &Entity{name: "player", hp: 80, deck: starterDeck()}
	enemy := &Entity{name: "enemy", hp: 40}
	globals.Define("player", player)
	globals.Define("enemy", enemy)
	return [2]*Entity{player, enemy}
}

func (e *Evaluator) combatTarget(value any) *Entity {
	target, ok := value.(*Entity)
	if !ok {
		panic(&RuntimeError{Token: e.callSite, Message: "Target must be the player or the enemy."})
	}
	return target
}

func (e *Evaluator) combatWholeNumber(value any, label string) float64 {
	amount, ok := value.(float64)
	if !ok || amount < 0 || math.IsInf(amount, 0) || amount != math.Trunc(amount) {
		panic(&RuntimeError{Token: e.callSite, Message: label + " must be a whole number, zero or more."})
	}
	return amount
}

func (e *Evaluator) combatArgs(args []any) (*Entity, float64) {
	return e.combatTarget(args[0]), e.combatWholeNumber(args[1], "Amount")
}

func init() {
	// deal(target, n) -> number: block soaks damage first, the rest comes off HP.
	// Returns the HP actually lost, so dealing 100 to a 34-HP target returns 34.
	natives["deal"] = &nativeFn{arity: 2, fn: func(e *Evaluator, args []any) any {
		target, amount := e.combatArgs(args)
		absorbed := math.Min(target.block, amount)
		lost := math.Min(target.hp, amount-absorbed)
		target.block -= absorbed
		target.hp -= lost
		return lost
	}}

	// block(target, n) -> number: adds block and returns the block now held.
	natives["block"] = &nativeFn{arity: 2, fn: func(e *Evaluator, args []any) any {
		target, amount := e.combatArgs(args)
		total := target.block + amount
		if math.IsInf(total, 0) {
			panic(&RuntimeError{Token: e.callSite, Message: "Block total is too large."})
		}
		target.block = total
		return total
	}}
}
