package comparison

import (
	"fmt"

	jsonfilter "github.com/youruser/json-filter"
)

// Instantiate creates a comparison operator implementation for the provided type.
func Instantiate(t Type, field string, value interface{}) (jsonfilter.Operator, error) {
	switch t {
	case Equal:
		op, err := NewEqualOperator(field, value)
		if err != nil {
			return nil, err
		}
		return op, nil
	case Regex:
		pattern, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("regex operator expects string value, got %T", value)
		}
		op, err := NewRegexOperator(field, pattern)
		if err != nil {
			return nil, err
		}
		return op, nil
	default:
		return nil, fmt.Errorf("comparison operator %s is not implemented", t)
	}
}
