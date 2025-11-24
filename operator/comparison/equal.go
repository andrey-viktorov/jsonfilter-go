package comparison

import (
	"fmt"
	"reflect"

	"github.com/tidwall/gjson"
	jsonfilter "github.com/youruser/json-filter"
)

// EqualOperator compares a JSON path value for equality against an expected literal.
type EqualOperator struct {
	jsonPath        string
	expected        interface{}
	pathNotFoundMsg string
	mismatchMsg     string
}

// NewEqualOperator constructs an EqualOperator instance.
func NewEqualOperator(jsonPath string, expected interface{}) (*EqualOperator, error) {
	if jsonPath == "" {
		return nil, fmt.Errorf("json path must not be empty")
	}
	op := &EqualOperator{
		jsonPath:        jsonPath,
		expected:        expected,
		pathNotFoundMsg: "json path " + jsonPath + " not found",
	}
	op.mismatchMsg = fmt.Sprintf("value did not equal expected %v", expected)
	return op, nil
}

// MustNewEqualOperator panics when inputs are invalid.
func MustNewEqualOperator(jsonPath string, expected interface{}) *EqualOperator {
	op, err := NewEqualOperator(jsonPath, expected)
	if err != nil {
		panic(err)
	}
	return op
}

// Name returns the operator identifier.
func (o *EqualOperator) Name() string {
	return string(Equal)
}

// Evaluate fetches the JSON value and compares it to the expected value.
func (o *EqualOperator) Evaluate(json []byte) jsonfilter.EvaluationResult {
	actual := getJSONResult(json, o.jsonPath)
	if !actual.Exists() {
		return jsonfilter.ErrorResult(o.Name(), o.pathNotFoundMsg)
	}

	if o.matches(actual) {
		return jsonfilter.ValidResult(o.Name())
	}

	return jsonfilter.ErrorResult(o.Name(), o.mismatchMsg)
}

// Validate ensures the operator is correctly configured.
func (o *EqualOperator) Validate() jsonfilter.ValidationResult {
	if o.jsonPath == "" {
		return jsonfilter.ErrorValidationResult(o.Name(), "json path must not be empty")
	}
	return jsonfilter.ValidValidationResult(o.Name())
}

func (o *EqualOperator) matches(actual gjson.Result) bool {
	switch expected := o.expected.(type) {
	case string:
		return actual.Str == expected
	case fmt.Stringer:
		return actual.Str == expected.String()
	case bool:
		return actual.Bool() == expected
	case int:
		return actual.Int() == int64(expected)
	case int8:
		return actual.Int() == int64(expected)
	case int16:
		return actual.Int() == int64(expected)
	case int32:
		return actual.Int() == int64(expected)
	case int64:
		return actual.Int() == expected
	case uint:
		return actual.Uint() == uint64(expected)
	case uint8:
		return actual.Uint() == uint64(expected)
	case uint16:
		return actual.Uint() == uint64(expected)
	case uint32:
		return actual.Uint() == uint64(expected)
	case uint64:
		return actual.Uint() == expected
	case float32:
		return actual.Float() == float64(expected)
	case float64:
		return actual.Float() == expected
	case nil:
		return !actual.Exists() || actual.Type == gjson.Null
	default:
		return reflect.DeepEqual(actual.Value(), expected)
	}
}
