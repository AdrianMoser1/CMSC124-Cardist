package evaluator

import (
	"math"
	"strconv"
)

// FormatValue is the single place where a Cardist runtime value becomes text.
// Both `./run --eval` and the REPL go through it, so the two can never drift.
//
// Value representation: float64 | string | bool | nil
//
// for notes, printing rules:
//   - numbers: shortest exact form, no trailing ".0"  	(5.0 -> "5", 1.5 -> "1.5")
//   - strings: printed raw, without quotes             ("hi" -> hi)
//   - booleans: true or false
//   - nil: "nil"
func FormatValue(v any) string {
	switch val := v.(type) {
	case nil:
		return "nil"
	case bool:
		return strconv.FormatBool(val)
	case string:
		return val
	case float64:
		return formatNumber(val)
	default:
		// Unreachable if the evaluator only produces the four types above.
		// Returning text (instead of panicking) keeps a host-side bug from
		// turning into a stack trace.
		return "<unknown value>"
	}
}

func formatNumber(f float64) string {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "-Infinity"
	case f == 0:
		// Collapses -0 to "0" so `-0` and `0` print identically.
		return "0"
	}
	// 'f' with precision -1 = smallest digits that round-trip, never scientific
	// notation, and no trailing ".0" for whole numbers.
	return strconv.FormatFloat(f, 'f', -1, 64)
}
