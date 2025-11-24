package logic

import (
	"testing"

	jsonfilter "github.com/youruser/json-filter"
	"github.com/youruser/json-filter/operator/comparison"
)

func BenchmarkAndOperatorEvaluate(b *testing.B) {
	child1 := comparison.MustNewEqualOperator("foo", "bar")
	child2 := comparison.MustNewRegexOperator("baz", `^qux`)

	op, err := NewOperator(And, []jsonfilter.Operator{child1, child2})
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}

	payload := []byte(`{"foo":"bar","baz":"qux"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if res := op.Evaluate(payload); !res.Match {
			b.Fatalf("expected AND to match: %#v", res)
		}
	}
}

func BenchmarkOrOperatorEvaluate(b *testing.B) {
	match := comparison.MustNewEqualOperator("foo", "bar")
	miss := comparison.MustNewRegexOperator("baz", `^miss`)

	op, err := NewOperator(Or, []jsonfilter.Operator{match, miss})
	if err != nil {
		b.Fatalf("unexpected error: %v", err)
	}

	payload := []byte(`{"foo":"bar","baz":"noop"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if res := op.Evaluate(payload); !res.Match {
			b.Fatalf("expected OR to match: %#v", res)
		}
	}
}
