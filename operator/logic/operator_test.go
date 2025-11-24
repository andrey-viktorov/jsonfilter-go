package logic

import (
	"testing"

	jsonfilter "github.com/youruser/json-filter"
)

type stubOperator struct {
	name       string
	evalResult jsonfilter.EvaluationResult
	valResult  jsonfilter.ValidationResult
	calls      *int
}

func (s *stubOperator) Name() string { return s.name }

func (s *stubOperator) Evaluate(_ []byte) jsonfilter.EvaluationResult {
	if s.calls != nil {
		*s.calls++
	}
	return s.evalResult
}

func (s *stubOperator) Validate() jsonfilter.ValidationResult {
	if s.valResult.OperatorName == "" {
		return jsonfilter.ValidValidationResult(s.name)
	}
	return s.valResult
}

func TestAndOperatorEvaluation(t *testing.T) {
	child1 := &stubOperator{name: "child1", evalResult: jsonfilter.ValidResult("child1")}
	child2 := &stubOperator{name: "child2", evalResult: jsonfilter.ErrorResult("child2", "boom")}

	op, err := NewOperator(And, []jsonfilter.Operator{child1, child2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := op.Evaluate([]byte(`{}`))
	if res.Match {
		t.Fatalf("expected AND operator to fail: %#v", res)
	}
	if res.CauseDescription == "" {
		t.Fatalf("expected cause description to propagate")
	}
}

func TestOrOperatorShortCircuit(t *testing.T) {
	callCount := 0
	matching := &stubOperator{
		name:       "match",
		evalResult: jsonfilter.ValidResult("match"),
		calls:      &callCount,
	}
	nonMatching := &stubOperator{
		name:       "miss",
		evalResult: jsonfilter.ErrorResult("miss", "nope"),
		calls:      &callCount,
	}

	op, err := NewOperator(Or, []jsonfilter.Operator{matching, nonMatching})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := op.Evaluate([]byte(`{}`))
	if !res.Match {
		t.Fatalf("expected OR operator to succeed: %#v", res)
	}
	if callCount != 1 {
		t.Fatalf("expected short circuit to stop after first match, got %d calls", callCount)
	}
}
