package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Op represents a comparison operator.
type Op string

const (
	OpEq  Op = "="
	OpNeq Op = "!="
	OpGt  Op = ">"
	OpLt  Op = "<"
	OpGte Op = ">="
	OpLte Op = "<="
	OpContains Op = "~"
)

// Condition represents a single filter condition.
type Condition struct {
	Field string
	Op    Op
	Value string
}

// Filter holds a set of conditions applied with AND logic.
type Filter struct {
	Conditions []Condition
}

// NewFilter creates a Filter from a slice of raw condition strings.
// Each condition is in the form "field=value", "field>value", etc.
func NewFilter(exprs []string) (*Filter, error) {
	f := &Filter{}
	for _, expr := range exprs {
		cond, err := parseCondition(expr)
		if err != nil {
			return nil, err
		}
		f.Conditions = append(f.Conditions, cond)
	}
	return f, nil
}

// Match returns true if the log entry satisfies all conditions.
func (f *Filter) Match(entry *parser.LogEntry) bool {
	for _, cond := range f.Conditions {
		if !matchCondition(entry, cond) {
			return false
		}
	}
	return true
}

func matchCondition(entry *parser.LogEntry, cond Condition) bool {
	val := entry.GetField(cond.Field)
	if val == nil {
		return false
	}
	actual := fmt.Sprintf("%v", val)

	if cond.Op == OpContains {
		return strings.Contains(strings.ToLower(actual), strings.ToLower(cond.Value))
	}

	// Try numeric comparison first.
	actualNum, errA := strconv.ParseFloat(actual, 64)
	expectedNum, errE := strconv.ParseFloat(cond.Value, 64)
	if errA == nil && errE == nil {
		return compareNums(actualNum, expectedNum, cond.Op)
	}

	return compareStrings(actual, cond.Value, cond.Op)
}

func compareNums(a, b float64, op Op) bool {
	switch op {
	case OpEq:
		return a == b
	case OpNeq:
		return a != b
	case OpGt:
		return a > b
	case OpLt:
		return a < b
	case OpGte:
		return a >= b
	case OpLte:
		return a <= b
	}
	return false
}

func compareStrings(a, b string, op Op) bool {
	switch op {
	case OpEq:
		return a == b
	case OpNeq:
		return a != b
	case OpGt:
		return a > b
	case OpLt:
		return a < b
	case OpGte:
		return a >= b
	case OpLte:
		return a <= b
	}
	return false
}
