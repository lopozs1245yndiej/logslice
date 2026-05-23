// Package output handles formatting and writing of log entries.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Format represents the output format for log entries.
type Format string

const (
	FormatJSON    Format = "json"
	FormatText    Format = "text"
	FormatCompact Format = "compact"
)

// Formatter writes log entries to an output stream.
type Formatter struct {
	writer    io.Writer
	format    Format
	timeField string
}

// NewFormatter creates a new Formatter with the given writer and format.
func NewFormatter(w io.Writer, format Format) *Formatter {
	return &Formatter{
		writer:    w,
		format:    format,
		timeField: "time",
	}
}

// SetTimeField configures which field is treated as the timestamp.
func (f *Formatter) SetTimeField(field string) {
	f.timeField = field
}

// Write formats and writes a single log entry.
func (f *Formatter) Write(entry map[string]interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(entry)
	case FormatText:
		return f.writeText(entry)
	case FormatCompact:
		return f.writeCompact(entry)
	default:
		return f.writeJSON(entry)
	}
}

// WriteAll formats and writes multiple log entries, stopping on first error.
func (f *Formatter) WriteAll(entries []map[string]interface{}) error {
	for _, entry := range entries {
		if err := f.Write(entry); err != nil {
			return fmt.Errorf("output: failed to write entry: %w", err)
		}
	}
	return nil
}

func (f *Formatter) writeJSON(entry map[string]interface{}) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintf(f.writer, "%s\n", b)
	return err
}

func (f *Formatter) writeText(entry map[string]interface{}) error {
	ts := extractTime(entry, f.timeField)
	level := extractString(entry, "level", "INFO")
	msg := extractString(entry, "message", extractString(entry, "msg", ""))

	var extras []string
	for k, v := range entry {
		if k == f.timeField || k == "level" || k == "message" || k == "msg" {
			continue
		}
		extras = append(extras, fmt.Sprintf("%s=%v", k, v))
	}

	line := fmt.Sprintf("%s [%s] %s", ts, strings.ToUpper(level), msg)
	if len(extras) > 0 {
		line += " " + strings.Join(extras, " ")
	}
	_, err := fmt.Fprintln(f.writer, line)
	return err
}

func (f *Formatter) writeCompact(entry map[string]interface{}) error {
	level := extractString(entry, "level", "?")
	msg := extractString(entry, "message", extractString(entry, "msg", ""))
	_, err := fmt.Fprintf(f.writer, "%-5s %s\n", strings.ToUpper(level), msg)
	return err
}

func extractString(entry map[string]interface{}, key, fallback string) string {
	if v, ok := entry[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return fallback
}

func extractTime(entry map[string]interface{}, key string) string {
	if v, ok := entry[key]; ok {
		switch t := v.(type) {
		case string:
			if parsed, err := time.Parse(time.RFC3339, t); err == nil {
				return parsed.Format("2006-01-02 15:04:05")
			}
			return t
		}
	}
	return time.Now().Format("2006-01-02 15:04:05")
}
