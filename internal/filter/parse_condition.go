package filter

import (
	"fmt"
	"strings"
)

// operators in descending length order so multi-char ops are checked first.
var operators = []Op{OpGte, OpLte, OpNeq, OpEq, OpGt, OpLt, OpContains}

// parseCondition parses a raw expression string into a Condition.
func parseCondition(expr string) (Condition, error) {
	for _, op := range operators {
		idx := strings.Index(expr, string(op))
		if idx < 0 {
			continue
		}
		field := strings.TrimSpace(expr[:idx])
		value := strings.TrimSpace(expr[idx+len(op):])
		if field == "" {
			return Condition{}, fmt.Errorf("filter: empty field name in expression %q", expr)
		}
		return Condition{Field: field, Op: op, Value: value}, nil
	}
	return Condition{}, fmt.Errorf("filter: no valid operator found in expression %q", expr)
}
