package filter

import (
	"testing"
)

func TestParseCondition_AllOperators(t *testing.T) {
	cases := []struct {
		expr  string
		field string
		op    Op
		value string
	}{
		{"level=error", "level", OpEq, "error"},
		{"level!=debug", "level", OpNeq, "debug"},
		{"status>400", "status", OpGt, "400"},
		{"status<500", "status", OpLt, "500"},
		{"status>=200", "status", OpGte, "200"},
		{"status<=299", "status", OpLte, "299"},
		{"message~error", "message", OpContains, "error"},
	}

	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			cond, err := parseCondition(tc.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cond.Field != tc.field {
				t.Errorf("field: got %q, want %q", cond.Field, tc.field)
			}
			if cond.Op != tc.op {
				t.Errorf("op: got %q, want %q", cond.Op, tc.op)
			}
			if cond.Value != tc.value {
				t.Errorf("value: got %q, want %q", cond.Value, tc.value)
			}
		})
	}
}

func TestParseCondition_EmptyField(t *testing.T) {
	_, err := parseCondition("=value")
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestParseCondition_NoOperator(t *testing.T) {
	_, err := parseCondition("justfield")
	if err == nil {
		t.Error("expected error when no operator present")
	}
}
