package jsonfilter

// Operator represents a single executable filter node.
type Operator interface {
	// Name returns the stable identifier of the operator.
	Name() string
	// Evaluate runs the operator against the provided JSON payload.
	Evaluate(json []byte) EvaluationResult
	// Validate checks the operator configuration prior to execution.
	Validate() ValidationResult
}
