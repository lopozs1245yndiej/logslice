package output

import (
	"bytes"
	"strings"
	"testing"
)

func makeEntry(fields map[string]interface{}) map[string]interface{} {
	return fields
}

func TestFormatter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	entry := makeEntry(map[string]interface{}{"level": "info", "msg": "hello"})
	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "\"level\":") {
		t.Errorf("expected JSON output, got: %s", out)
	}
	if !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON to end with }, got: %s", out)
	}
}

func TestFormatter_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatText)

	entry := makeEntry(map[string]interface{}{
		"time":    "2024-01-15T10:00:00Z",
		"level":   "error",
		"message": "something failed",
		"code":    500,
	})
	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in text output, got: %s", out)
	}
	if !strings.Contains(out, "something failed") {
		t.Errorf("expected message in text output, got: %s", out)
	}
}

func TestFormatter_CompactFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatCompact)

	entry := makeEntry(map[string]interface{}{"level": "warn", "msg": "disk low"})
	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN in compact output, got: %s", out)
	}
	if !strings.Contains(out, "disk low") {
		t.Errorf("expected message in compact output, got: %s", out)
	}
}

func TestFormatter_SetTimeField(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatText)
	f.SetTimeField("@timestamp")

	entry := makeEntry(map[string]interface{}{
		"@timestamp": "2024-06-01T08:30:00Z",
		"level":      "info",
		"msg":        "custom time field",
	})
	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "2024-06-01") {
		t.Errorf("expected parsed timestamp in output, got: %s", out)
	}
}

func TestFormatter_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, Format("unknown"))

	entry := makeEntry(map[string]interface{}{"level": "debug", "msg": "test"})
	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "{") {
		t.Errorf("expected JSON fallback output, got: %s", out)
	}
}
