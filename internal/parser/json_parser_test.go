package parser

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestParser_SingleEntry(t *testing.T) {
	input := `{"time":"2024-01-15T10:00:00Z","level":"info","message":"server started","port":8080}`
	p := NewParser(strings.NewReader(input))

	entry, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entry.Level != "info" {
		t.Errorf("expected level 'info', got %q", entry.Level)
	}
	if entry.Message != "server started" {
		t.Errorf("expected message 'server started', got %q", entry.Message)
	}
	expected := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, entry.Timestamp)
	}
}

func TestParser_MultipleEntries(t *testing.T) {
	input := `{"level":"info","msg":"start"}
{"level":"error","msg":"failed"}
`
	p := NewParser(strings.NewReader(input))

	levels := []string{"info", "error"}
	for i, want := range levels {
		entry, err := p.Next()
		if err != nil {
			t.Fatalf("entry %d: unexpected error: %v", i, err)
		}
		if entry.Level != want {
			t.Errorf("entry %d: expected level %q, got %q", i, want, entry.Level)
		}
	}

	_, err := p.Next()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestParser_AlternativeFieldNames(t *testing.T) {
	input := `{"severity":"warn","message":"disk usage high","ts":"2024-06-01T12:00:00Z"}`
	p := NewParser(strings.NewReader(input))

	entry, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != "warn" {
		t.Errorf("expected level 'warn', got %q", entry.Level)
	}
}

func TestLogEntry_GetField(t *testing.T) {
	input := `{"level":"debug","msg":"trace","request_id":"abc-123"}`
	p := NewParser(strings.NewReader(input))

	entry, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, ok := entry.GetString("request_id")
	if !ok {
		t.Fatal("expected request_id field to exist")
	}
	if val != "abc-123" {
		t.Errorf("expected 'abc-123', got %q", val)
	}

	_, ok = entry.GetField("nonexistent")
	if ok {
		t.Error("expected nonexistent field to be absent")
	}
}

func TestParser_InvalidJSON(t *testing.T) {
	input := `not valid json`
	p := NewParser(strings.NewReader(input))

	_, err := p.Next()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
