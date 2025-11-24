package serde

import "testing"

func TestParserFromJSON(t *testing.T) {
	parser := DefaultParser()
	payload := []byte(`{"and":[{"eq":{"field":"foo","value":"bar"}},{"rx":{"field":"baz","value":"^qux"}}]}`)

	op, err := parser.FromJSON(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := op.Evaluate([]byte(`{"foo":"bar","baz":"quxa"}`))
	if !res.Match {
		t.Fatalf("expected match, got %#v", res)
	}
}

func TestParserRespectsComplexityLimit(t *testing.T) {
	parser := NewParser(1)
	payload := []byte(`{"and":[{"eq":{"field":"foo","value":"bar"}},{"eq":{"field":"bar","value":"baz"}}]}`)

	if _, err := parser.FromJSON(payload); err == nil {
		t.Fatalf("expected complexity error")
	}
}

func TestParserFromYAML(t *testing.T) {
	parser := DefaultParser()
	payload := []byte(`
jsonFilter:
  or:
    - eq:
        field: foo
        value: nope
    - eq:
        field: foo
        value: ok
`)

	op, err := parser.FromYAML(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := op.Evaluate([]byte(`{"foo":"ok"}`))
	if !res.Match {
		t.Fatalf("expected match for YAML parsed operator: %#v", res)
	}
}

func TestParserFromMap(t *testing.T) {
	parser := DefaultParser()
	root := map[string]interface{}{
		"and": []interface{}{
			map[string]interface{}{
				"eq": map[string]interface{}{
					"field": "foo",
					"value": "bar",
				},
			},
			map[string]interface{}{
				"rx": map[string]interface{}{
					"field": "baz",
					"value": "^qu",
				},
			},
		},
	}

	op, err := parser.FromMap(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := op.Evaluate([]byte(`{"foo":"bar","baz":"quasar"}`))
	if !res.Match {
		t.Fatalf("expected map-parsed operator to match: %#v", res)
	}
}
