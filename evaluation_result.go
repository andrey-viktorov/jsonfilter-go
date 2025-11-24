package jsonfilter

// EvaluationResult captures the outcome of running an Operator against a JSON payload.
type EvaluationResult struct {
	Match            bool               `json:"match" yaml:"match"`
	OperatorName     string             `json:"operatorName" yaml:"operatorName"`
	CauseDescription string             `json:"causeDescription,omitempty" yaml:"causeDescription,omitempty"`
	ChildOperators   []EvaluationResult `json:"childOperators,omitempty" yaml:"childOperators,omitempty"`
}

// ValidResult returns a successful EvaluationResult for the supplied operator name.
func ValidResult(operatorName string) EvaluationResult {
	return EvaluationResult{Match: true, OperatorName: operatorName}
}

// ErrorResult returns an EvaluationResult describing why an operator failed.
func ErrorResult(operatorName, cause string) EvaluationResult {
	return EvaluationResult{Match: false, OperatorName: operatorName, CauseDescription: cause}
}

// AggregateResult aggregates child operator results using a precomputed match value.
func AggregateResult(operatorName string, match bool, children []EvaluationResult, cause string) EvaluationResult {
	return EvaluationResult{
		Match:            match,
		OperatorName:     operatorName,
		CauseDescription: cause,
		ChildOperators:   children,
	}
}
