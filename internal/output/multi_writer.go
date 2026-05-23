package output

import (
	"fmt"
	"io"
	"sync"
)

// MultiWriter fans out log entries to multiple Formatters concurrently-safe.
type MultiWriter struct {
	mu         sync.Mutex
	formatters []*Formatter
}

// NewMultiWriter creates a MultiWriter with no formatters attached.
func NewMultiWriter() *MultiWriter {
	return &MultiWriter{}
}

// Add registers a new Formatter backed by the given writer and format.
func (mw *MultiWriter) Add(w io.Writer, format Format) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	mw.formatters = append(mw.formatters, NewFormatter(w, format))
}

// Write sends the entry to all registered formatters.
// It collects all errors and returns a combined error if any occurred.
func (mw *MultiWriter) Write(entry map[string]interface{}) error {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	var errs []error
	for _, f := range mw.formatters {
		if err := f.Write(entry); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("output: %d formatter(s) failed: %v", len(errs), errs)
}

// Len returns the number of registered formatters.
func (mw *MultiWriter) Len() int {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	return len(mw.formatters)
}
