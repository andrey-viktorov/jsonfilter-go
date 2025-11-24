package comparison

import "testing"

func TestEqualOperatorMatch(t *testing.T) {
	op := MustNewEqualOperator("foo", "bar")
	result := op.Evaluate([]byte(`{"foo":"bar"}`))
	if !result.Match {
		t.Fatalf("expected match, got %#v", result)
	}
}

func TestEqualOperatorMismatch(t *testing.T) {
	op := MustNewEqualOperator("foo", "bar")
	result := op.Evaluate([]byte(`{"foo":"baz"}`))
	if result.Match {
		t.Fatalf("expected mismatch, got %#v", result)
	}
}

func TestRegexOperatorEvaluate(t *testing.T) {
	op := MustNewRegexOperator("foo", `^ba.`)
	if res := op.Evaluate([]byte(`{"foo":"baz"}`)); !res.Match {
		t.Fatalf("expected regex to match: %#v", res)
	}
	if res := op.Evaluate([]byte(`{"foo":"qux"}`)); res.Match {
		t.Fatalf("expected regex to fail: %#v", res)
	}
}
