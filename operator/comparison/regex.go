package comparison

import (
	"fmt"
	"regexp"

	jsonfilter "github.com/youruser/json-filter"
)

// RegexOperator evaluates the value of a JSON path against a compiled regular expression.
type RegexOperator struct {
	jsonPath           string
	pattern            string
	compiledRe         *regexp.Regexp
	pathNotFoundMsg    string
	patternMismatchMsg string
}

// NewRegexOperator creates a RegexOperator and compiles the provided pattern.
func NewRegexOperator(jsonPath, pattern string) (*RegexOperator, error) {
	if jsonPath == "" {
		return nil, fmt.Errorf("json path must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("regex pattern must not be empty")
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}
	op := &RegexOperator{
		jsonPath:        jsonPath,
		pattern:         pattern,
		compiledRe:      compiled,
		pathNotFoundMsg: "json path " + jsonPath + " not found",
	}
	op.patternMismatchMsg = fmt.Sprintf("value does not match regex %s", pattern)
	return op, nil
}

// MustNewRegexOperator panics if construction fails.
func MustNewRegexOperator(jsonPath, pattern string) *RegexOperator {
	op, err := NewRegexOperator(jsonPath, pattern)
	if err != nil {
		panic(err)
	}
	return op
}

// Name returns the operator identifier.
func (o *RegexOperator) Name() string {
	return string(Regex)
}

// Evaluate executes the regex match against the JSON value at jsonPath.
func (o *RegexOperator) Evaluate(json []byte) jsonfilter.EvaluationResult {
	actual := getJSONResult(json, o.jsonPath)
	if !actual.Exists() {
		return jsonfilter.ErrorResult(o.Name(), o.pathNotFoundMsg)
	}

	if o.compiledRe.MatchString(actual.Str) {
		return jsonfilter.ValidResult(o.Name())
	}

	return jsonfilter.ErrorResult(o.Name(), o.patternMismatchMsg)
}

// Validate re-validates invariant fields.
func (o *RegexOperator) Validate() jsonfilter.ValidationResult {
	if o.jsonPath == "" {
		return jsonfilter.ErrorValidationResult(o.Name(), "json path must not be empty")
	}
	if o.pattern == "" || o.compiledRe == nil {
		return jsonfilter.ErrorValidationResult(o.Name(), "regex operator must have a compiled pattern")
	}
	return jsonfilter.ValidValidationResult(o.Name())
}
