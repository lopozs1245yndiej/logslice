// Package output provides formatting and writing utilities for structured log entries.
//
// It supports three output formats:
//
//   - json: Each log entry is written as a compact JSON object (default).
//   - text: Human-readable format: "TIMESTAMP [LEVEL] message key=value ..."
//   - compact: Minimal format showing only level and message.
//
// Usage:
//
//	w := output.NewFormatter(os.Stdout, output.FormatText)
//	w.Write(entry)
//
// The time field used for timestamp extraction defaults to "time" but can
// be overridden via SetTimeField to support alternative schemas such as
// Elasticsearch's "@timestamp".
package output
