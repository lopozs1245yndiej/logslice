package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// LogEntry represents a single parsed JSON log line.
type LogEntry struct {
	Raw       map[string]interface{}
	Timestamp time.Time
	Level     string
	Message   string
}

// Parser reads and parses newline-delimited JSON log entries.
type Parser struct {
	decoder *json.Decoder
}

// NewParser creates a new Parser reading from the given reader.
func NewParser(r io.Reader) *Parser {
	return &Parser{
		decoder: json.NewDecoder(r),
	}
}

// Next reads the next log entry from the stream.
// Returns io.EOF when the stream is exhausted.
func (p *Parser) Next() (*LogEntry, error) {
	var raw map[string]interface{}
	if err := p.decoder.Decode(&raw); err != nil {
		return nil, err
	}

	entry := &LogEntry{Raw: raw}

	if ts, ok := extractString(raw, "time", "timestamp", "ts"); ok {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			entry.Timestamp = t
		}
	}

	if level, ok := extractString(raw, "level", "severity", "lvl"); ok {
		entry.Level = level
	}

	if msg, ok := extractString(raw, "message", "msg"); ok {
		entry.Message = msg
	}

	return entry, nil
}

// GetField retrieves a field value from the log entry by key.
func (e *LogEntry) GetField(key string) (interface{}, bool) {
	v, ok := e.Raw[key]
	return v, ok
}

// GetString retrieves a string field from the log entry.
func (e *LogEntry) GetString(key string) (string, bool) {
	v, ok := e.Raw[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// String returns a compact JSON representation of the entry.
func (e *LogEntry) String() string {
	b, err := json.Marshal(e.Raw)
	if err != nil {
		return fmt.Sprintf("%v", e.Raw)
	}
	return string(b)
}

// extractString tries each candidate key and returns the first string value found.
func extractString(m map[string]interface{}, keys ...string) (string, bool) {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
	}
	return "", false
}
