package jsonfilter

// ValidationResult captures pre-execution validation feedback for Operator trees.
type ValidationResult struct {
	Valid            bool               `json:"valid" yaml:"valid"`
	OperatorName     string             `json:"operatorName" yaml:"operatorName"`
	CauseDescription string             `json:"causeDescription,omitempty" yaml:"causeDescription,omitempty"`
	ChildOperators   []ValidationResult `json:"childOperators,omitempty" yaml:"childOperators,omitempty"`
}

// ValidValidationResult constructs a successful ValidationResult.
func ValidValidationResult(operatorName string) ValidationResult {
	return ValidationResult{Valid: true, OperatorName: operatorName}
}

// ErrorValidationResult constructs a failed ValidationResult with the provided cause.
func ErrorValidationResult(operatorName, cause string) ValidationResult {
	return ValidationResult{Valid: false, OperatorName: operatorName, CauseDescription: cause}
}

// AggregateValidationResult builds a ValidationResult aggregating child validations.
func AggregateValidationResult(operatorName string, valid bool, children []ValidationResult, cause string) ValidationResult {
	return ValidationResult{
		Valid:            valid,
		OperatorName:     operatorName,
		CauseDescription: cause,
		ChildOperators:   children,
	}
}
