// Package source provides abstractions for reading log entries
// from various input sources such as files, stdin, or network streams.
package source

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Source represents a named log input source.
type Source struct {
	Name   string
	reader io.ReadCloser
	scanner *bufio.Scanner
}

// NewFileSource opens a file at the given path and returns a Source.
func NewFileSource(path string) (*Source, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("source: open file %q: %w", path, err)
	}
	return newSource(path, f), nil
}

// NewReaderSource creates a Source from an arbitrary io.ReadCloser.
func NewReaderSource(name string, r io.ReadCloser) *Source {
	return newSource(name, r)
}

// NewStdinSource creates a Source that reads from os.Stdin.
func NewStdinSource() *Source {
	return newSource("stdin", io.NopCloser(os.Stdin))
}

func newSource(name string, r io.ReadCloser) *Source {
	return &Source{
		Name:    name,
		reader:  r,
		scanner: bufio.NewScanner(r),
	}
}

// Lines returns a channel that emits raw log lines from the source.
// The channel is closed when the source is exhausted or an error occurs.
func (s *Source) Lines() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for s.scanner.Scan() {
			line := s.scanner.Text()
			if line != "" {
				ch <- line
			}
		}
	}()
	return ch
}

// Close releases the underlying reader.
func (s *Source) Close() error {
	return s.reader.Close()
}
