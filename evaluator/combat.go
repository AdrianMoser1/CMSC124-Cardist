package evaluator

import (
	"fmt"
	"math"
)

// Entity is something that can take damage and hold block. The player also
// owns the piles draw() moves cards between.
//
// shortcut: piles are plain counts, not cards, and the starting stats are
// fixed; upgrade both when `card` and `enemy` definition blocks exist.
// shortcut: block never expires; upgrade when `turn` blocks exist.
type Entity struct {
	name       string
	hp, block  float64
	deck, hand float64
}

func (en *Entity) String() string {
	return fmt.Sprintf("<%s hp=%s block=%s>", en.name, formatNumber(en.hp), formatNumber(en.block))
}

// defineCombatants makes `player` and `enemy` ordinary globals. The parser
// accepts those two reserved words as names, so scripts can use them.
func defineCombatants(globals *Environment) {
	globals.Define("player", &Entity{name: "player", hp: 80, deck: 10})
	globals.Define("enemy", &Entity{name: "enemy", hp: 40})
}

// combatArgs validates the (target, amount) pair every combat native takes.
func (e *Evaluator) combatArgs(args []any) (*Entity, float64) {
	target, ok := args[0].(*Entity)
	if !ok {
		panic(&RuntimeError{Token: e.callSite, Message: "Target must be the player or the enemy."})
	}
	amount, ok := args[1].(float64)
	if !ok || amount < 0 || amount != math.Trunc(amount) { // NaN fails the last test too
		panic(&RuntimeError{Token: e.callSite, Message: "Amount must be a whole number, zero or more."})
	}
	return target, amount
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
		target.block += amount
		return target.block
	}}

	// draw(target, n) -> number: moves up to n cards from draw pile to hand.
	// Returns how many moved, which is fewer than n when the pile runs short.
	natives["draw"] = &nativeFn{arity: 2, fn: func(e *Evaluator, args []any) any {
		target, amount := e.combatArgs(args)
		moved := math.Min(target.deck, amount)
		target.deck -= moved
		target.hand += moved
		return moved
	}}
}
