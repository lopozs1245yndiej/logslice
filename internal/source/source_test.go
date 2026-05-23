package source

import (
	"io"
	"sort"
	"strings"
	"testing"
)

func readerSource(name, data string) *Source {
	return NewReaderSource(name, io.NopCloser(strings.NewReader(data)))
}

func TestSource_Lines_Single(t *testing.T) {
	s := readerSource("test", "{\"level\":\"info\"}\n{\"level\":\"warn\"}\n")
	var got []string
	for line := range s.Lines() {
		got = append(got, line)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestSource_Lines_SkipsEmpty(t *testing.T) {
	s := readerSource("test", "line1\n\nline2\n\n")
	var got []string
	for line := range s.Lines() {
		got = append(got, line)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 non-empty lines, got %d", len(got))
	}
}

func TestSource_Name(t *testing.T) {
	s := readerSource("myfile.log", "")
	if s.Name != "myfile.log" {
		t.Errorf("expected name %q, got %q", "myfile.log", s.Name)
	}
}

func TestMultiSource_Len(t *testing.T) {
	m := NewMultiSource(readerSource("a", ""), readerSource("b", ""))
	if m.Len() != 2 {
		t.Errorf("expected len 2, got %d", m.Len())
	}
	m.Add(readerSource("c", ""))
	if m.Len() != 3 {
		t.Errorf("expected len 3, got %d", m.Len())
	}
}

func TestMultiSource_Lines_MergesAll(t *testing.T) {
	a := readerSource("a", "line-a1\nline-a2\n")
	b := readerSource("b", "line-b1\n")
	m := NewMultiSource(a, b)

	var lines []string
	for tl := range m.Lines() {
		lines = append(lines, tl.Source+":"+tl.Line)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 tagged lines, got %d: %v", len(lines), lines)
	}
	sort.Strings(lines)
	expected := []string{"a:line-a1", "a:line-a2", "b:line-b1"}
	for i, e := range expected {
		if lines[i] != e {
			t.Errorf("line[%d]: expected %q, got %q", i, e, lines[i])
		}
	}
}
