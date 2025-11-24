package logic

import (
	"fmt"

	jsonfilter "github.com/youruser/json-filter"
)

// Operator represents a logical aggregation of child operators (and/or).
type Operator struct {
	typ      Type
	children []jsonfilter.Operator
}

// NewOperator builds a new logic operator instance.
func NewOperator(opType Type, children []jsonfilter.Operator) (*Operator, error) {
	if _, ok := allTypes[opType]; !ok {
		return nil, fmt.Errorf("unsupported logic operator %q", opType)
	}
	copied := make([]jsonfilter.Operator, len(children))
	copy(copied, children)
	return &Operator{typ: opType, children: copied}, nil
}

// MustNewOperator panics when construction fails.
func MustNewOperator(opType Type, children []jsonfilter.Operator) *Operator {
	op, err := NewOperator(opType, children)
	if err != nil {
		panic(err)
	}
	return op
}

// Name returns the identifier of the logic operator.
func (o *Operator) Name() string {
	return string(o.typ)
}

// Evaluate executes the logic operator against the provided JSON payload.
func (o *Operator) Evaluate(json []byte) jsonfilter.EvaluationResult {
	if len(o.children) == 0 {
		return jsonfilter.ErrorResult(o.Name(), "logic operator requires at least one child")
	}

	switch o.typ {
	case And:
		return o.evaluateAnd(json)
	case Or:
		return o.evaluateOr(json)
	default:
		return jsonfilter.ErrorResult(o.Name(), fmt.Sprintf("unsupported logic operator %q", o.typ))
	}
}

func (o *Operator) evaluateAnd(json []byte) jsonfilter.EvaluationResult {
	for _, child := range o.children {
		result := child.Evaluate(json)
		if !result.Match {
			cause := result.CauseDescription
			if cause == "" {
				cause = "child operator returned no match"
			}
			return jsonfilter.ErrorResult(o.Name(), cause)
		}
	}
	return jsonfilter.ValidResult(o.Name())
}

func (o *Operator) evaluateOr(json []byte) jsonfilter.EvaluationResult {
	for _, child := range o.children {
		result := child.Evaluate(json)
		if result.Match {
			return jsonfilter.ValidResult(o.Name())
		}
	}
	return jsonfilter.ErrorResult(o.Name(), "no child operator produced a match")
}

// Validate ensures the logic operator and child operators are well defined.
func (o *Operator) Validate() jsonfilter.ValidationResult {
	if len(o.children) == 0 {
		return jsonfilter.ErrorValidationResult(o.Name(), "logic operator requires at least one child")
	}

	switch o.typ {
	case And:
		return o.validateAnd()
	case Or:
		return o.validateOr()
	default:
		return jsonfilter.ErrorValidationResult(o.Name(), fmt.Sprintf("unsupported logic operator %q", o.typ))
	}
}

func (o *Operator) validateAnd() jsonfilter.ValidationResult {
	children := make([]jsonfilter.ValidationResult, 0, len(o.children))
	allValid := true
	for _, child := range o.children {
		result := child.Validate()
		children = append(children, result)
		if !result.Valid {
			allValid = false
		}
	}
	cause := ""
	if !allValid {
		cause = "child operator validation failed"
	}
	return jsonfilter.AggregateValidationResult(o.Name(), allValid, children, cause)
}

func (o *Operator) validateOr() jsonfilter.ValidationResult {
	children := make([]jsonfilter.ValidationResult, 0, len(o.children))
	anyValid := false
	for _, child := range o.children {
		result := child.Validate()
		children = append(children, result)
		if result.Valid {
			anyValid = true
		}
	}
	cause := ""
	if !anyValid {
		cause = "no child operator validated successfully"
	}
	return jsonfilter.AggregateValidationResult(o.Name(), anyValid, children, cause)
}
