// Package filter provides query-condition parsing and log entry matching
// for logslice.
//
// A Filter is constructed from one or more expression strings, each of the
// form:
//
//	<field><op><value>
//
// Supported operators:
//
//	=    equal
//	!=   not equal
//	>    greater than
//	<    less than
//	>=   greater than or equal
//	<=   less than or equal
//	~    contains (case-insensitive substring match)
//
// Numeric fields are compared as float64 when both sides parse as numbers;
// otherwise a lexicographic string comparison is used.
//
// Multiple conditions are combined with AND semantics: all conditions must
// match for an entry to be included.
//
// Example:
//
//	f, err := filter.NewFilter([]string{"level=error", "status>=500"})
//	if err != nil { ... }
//	if f.Match(entry) {
//	    // entry satisfies both conditions
//	}
package filter
