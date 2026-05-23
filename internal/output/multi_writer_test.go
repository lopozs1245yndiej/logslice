package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestMultiWriter_WritesToAll(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	mw := NewMultiWriter()
	mw.Add(&buf1, FormatJSON)
	mw.Add(&buf2, FormatCompact)

	entry := map[string]interface{}{"level": "info", "msg": "broadcast"}
	if err := mw.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf1.String(), "broadcast") {
		t.Errorf("buf1 missing entry: %s", buf1.String())
	}
	if !strings.Contains(buf2.String(), "broadcast") {
		t.Errorf("buf2 missing entry: %s", buf2.String())
	}
}

func TestMultiWriter_Len(t *testing.T) {
	mw := NewMultiWriter()
	if mw.Len() != 0 {
		t.Errorf("expected 0 formatters, got %d", mw.Len())
	}

	var buf bytes.Buffer
	mw.Add(&buf, FormatJSON)
	mw.Add(&buf, FormatText)

	if mw.Len() != 2 {
		t.Errorf("expected 2 formatters, got %d", mw.Len())
	}
}

func TestMultiWriter_Empty(t *testing.T) {
	mw := NewMultiWriter()
	entry := map[string]interface{}{"level": "debug", "msg": "no-op"}
	if err := mw.Write(entry); err != nil {
		t.Errorf("expected no error with zero formatters, got: %v", err)
	}
}
