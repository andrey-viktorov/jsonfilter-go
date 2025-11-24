package comparison

import "fmt"

// Type enumerates supported comparison operators.
type Type string

const (
	Equal        Type = "eq"
	NotEqual     Type = "ne"
	Regex        Type = "rx"
	LessThan     Type = "lt"
	LessEqual    Type = "le"
	GreaterThan  Type = "gt"
	GreaterEqual Type = "ge"
	In           Type = "in"
	NotIn        Type = "nin"
	Contains     Type = "ct"
	NotContains  Type = "nct"
)

var allTypes = map[Type]struct{}{
	Equal:        {},
	NotEqual:     {},
	Regex:        {},
	LessThan:     {},
	LessEqual:    {},
	GreaterThan:  {},
	GreaterEqual: {},
	In:           {},
	NotIn:        {},
	Contains:     {},
	NotContains:  {},
}

// ParseType validates and returns the corresponding Type.
func ParseType(op string) (Type, error) {
	t := Type(op)
	if _, ok := allTypes[t]; !ok {
		return "", fmt.Errorf("comparison operator %q is not supported", op)
	}
	return t, nil
}

// MustParseType panics when the provided operator string is not supported.
func MustParseType(op string) Type {
	t, err := ParseType(op)
	if err != nil {
		panic(err)
	}
	return t
}
