package evaluator

import (
	"fmt"
	"math"
	"slices"
)

// Card represents one physical copy, so two cards with the same name can be banished separately.
type Card struct {
	name string
	cost float64
}

func (c *Card) String() string {
	return fmt.Sprintf("<card %s cost=%s>", c.name, formatNumber(c.cost))
}

func starterDeck() []*Card {
	deck := make([]*Card, 0, 10)
	for range 5 {
		deck = append(deck, &Card{name: "Strike", cost: 1})
	}
	for range 5 {
		deck = append(deck, &Card{name: "Defend", cost: 1})
	}
	return deck
}

func (e *Evaluator) combatCard(value any) *Card {
	card, ok := value.(*Card)
	if !ok {
		panic(&RuntimeError{Token: e.callSite, Message: "Expected a card."})
	}
	return card
}

func init() {
	natives["draw"] = &nativeFn{arity: 2, fn: func(e *Evaluator, args []any) any {
		target, amount := e.combatArgs(args)
		// Limit the float before converting it, since a valid amount can exceed Go's int range.
		moved := int(math.Min(float64(len(target.deck)), amount))
		target.hand = append(target.hand, target.deck[:moved]...)
		clear(target.deck[:moved])
		target.deck = target.deck[moved:]
		return float64(moved)
	}}

	natives["hand"] = &nativeFn{arity: 2, fn: func(e *Evaluator, args []any) any {
		target := e.combatTarget(args[0])
		index := e.combatWholeNumber(args[1], "Hand index")
		if index >= float64(len(target.hand)) {
			return nil
		}
		return target.hand[int(index)]
	}}

	natives["cost"] = &nativeFn{arity: 1, fn: func(e *Evaluator, args []any) any {
		return e.combatCard(args[0]).cost
	}}

	natives["banish"] = &nativeFn{arity: 2, fn: func(e *Evaluator, args []any) any {
		target := e.combatTarget(args[0])
		card := e.combatCard(args[1])
		index := slices.Index(target.hand, card)
		if index < 0 {
			return false
		}
		target.hand = slices.Delete(target.hand, index, index+1)
		return true
	}}
}
