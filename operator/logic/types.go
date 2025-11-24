package logic

import "fmt"

// Type enumerates supported logic operators.
type Type string

const (
	And Type = "and"
	Or  Type = "or"
)

var allTypes = map[Type]struct{}{
	And: {},
	Or:  {},
}

// ParseType validates the provided operator name.
func ParseType(op string) (Type, error) {
	t := Type(op)
	if _, ok := allTypes[t]; !ok {
		return "", fmt.Errorf("logic operator %q is not supported", op)
	}
	return t, nil
}

// MustParseType panics if the provided operator name is invalid.
func MustParseType(op string) Type {
	t, err := ParseType(op)
	if err != nil {
		panic(err)
	}
	return t
}
