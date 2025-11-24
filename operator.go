package jsonfilter

// Operator represents a single executable filter node.
type Operator interface {
	Name() string
	Evaluate(json []byte) EvaluationResult
	Validate() ValidationResult
}
