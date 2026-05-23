package source

import (
	"sync"
)

// TaggedLine pairs a raw log line with the name of its originating source.
type TaggedLine struct {
	Source string
	Line   string
}

// MultiSource fans in lines from multiple Sources into a single channel.
type MultiSource struct {
	sources []*Source
}

// NewMultiSource creates a MultiSource from the provided sources.
func NewMultiSource(sources ...*Source) *MultiSource {
	return &MultiSource{sources: sources}
}

// Add appends a source to the MultiSource.
func (m *MultiSource) Add(s *Source) {
	m.sources = append(m.sources, s)
}

// Len returns the number of sources.
func (m *MultiSource) Len() int {
	return len(m.sources)
}

// Lines merges lines from all sources concurrently into one channel.
// Each line is tagged with its source name. The channel closes when all
// sources are exhausted.
func (m *MultiSource) Lines() <-chan TaggedLine {
	out := make(chan TaggedLine)
	var wg sync.WaitGroup

	for _, s := range m.sources {
		wg.Add(1)
		go func(src *Source) {
			defer wg.Done()
			for line := range src.Lines() {
				out <- TaggedLine{Source: src.Name, Line: line}
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// CloseAll closes all underlying sources.
func (m *MultiSource) CloseAll() {
	for _, s := range m.sources {
		_ = s.Close()
	}
}
