package comparison

import "testing"

func BenchmarkEqualOperatorEvaluateMatch(b *testing.B) {
	op := MustNewEqualOperator("foo", "bar")
	payload := []byte(`{"foo":"bar","num":1}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if res := op.Evaluate(payload); !res.Match {
			b.Fatalf("expected match, got %#v", res)
		}
	}
}

func BenchmarkEqualOperatorEvaluateMiss(b *testing.B) {
	op := MustNewEqualOperator("foo", "bar")
	payload := []byte(`{"foo":"baz"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if res := op.Evaluate(payload); res.Match {
			b.Fatalf("expected mismatch, got %#v", res)
		}
	}
}

func BenchmarkRegexOperatorEvaluate(b *testing.B) {
	op := MustNewRegexOperator("foo", `^ba.*`)
	payload := []byte(`{"foo":"barbaz"}`)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if res := op.Evaluate(payload); !res.Match {
			b.Fatalf("expected regex match, got %#v", res)
		}
	}
}
