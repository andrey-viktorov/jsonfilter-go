// Package jsonfilter evaluates JSON payloads against operator trees parsed from
// YAML or JSON filter definitions. It exposes a minimal Operator interface that
// can be composed using comparison and logic operators and executed against raw
// []byte payloads with minimal allocations.
package jsonfilter
