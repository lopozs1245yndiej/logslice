package filter_test

import (
	"testing"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) *parser.LogEntry {
	return &parser.LogEntry{Fields: fields}
}

func TestFilter_EqualityMatch(t *testing.T) {
	f, err := filter.NewFilter([]string{"level=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !f.Match(makeEntry(map[string]interface{}{"level": "error"})) {
		t.Error("expected match for level=error")
	}
	if f.Match(makeEntry(map[string]interface{}{"level": "info"})) {
		t.Error("expected no match for level=info")
	}
}

func TestFilter_NumericGreaterThan(t *testing.T) {
	f, _ := filter.NewFilter([]string{"status>399"})

	if !f.Match(makeEntry(map[string]interface{}{"status": float64(500)})) {
		t.Error("expected match for status=500 > 399")
	}
	if f.Match(makeEntry(map[string]interface{}{"status": float64(200)})) {
		t.Error("expected no match for status=200 > 399")
	}
}

func TestFilter_ContainsOperator(t *testing.T) {
	f, _ := filter.NewFilter([]string{"message~timeout"})

	if !f.Match(makeEntry(map[string]interface{}{"message": "connection timeout occurred"})) {
		t.Error("expected match for message containing 'timeout'")
	}
	if f.Match(makeEntry(map[string]interface{}{"message": "all good"})) {
		t.Error("expected no match")
	}
}

func TestFilter_MultipleConditions(t *testing.T) {
	f, _ := filter.NewFilter([]string{"level=error", "status>=500"})

	if !f.Match(makeEntry(map[string]interface{}{"level": "error", "status": float64(503)})) {
		t.Error("expected match")
	}
	if f.Match(makeEntry(map[string]interface{}{"level": "error", "status": float64(200)})) {
		t.Error("expected no match: status too low")
	}
}

func TestFilter_MissingField(t *testing.T) {
	f, _ := filter.NewFilter([]string{"trace_id=abc"})
	if f.Match(makeEntry(map[string]interface{}{"level": "info"})) {
		t.Error("expected no match when field is missing")
	}
}

func TestNewFilter_InvalidExpression(t *testing.T) {
	_, err := filter.NewFilter([]string{"nodoperator"})
	if err == nil {
		t.Error("expected error for expression without operator")
	}
}
