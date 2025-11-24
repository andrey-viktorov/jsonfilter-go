package serde

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	jsonfilter "github.com/youruser/json-filter"
	"github.com/youruser/json-filter/operator/comparison"
	"github.com/youruser/json-filter/operator/logic"
	"gopkg.in/yaml.v3"
)

const defaultMaxComplexity = 42

// Parser turns YAML/JSON filter definitions into executable operator trees.
type Parser struct {
	maxComplexity int
}

// NewParser builds a parser enforcing the configured complexity limit.
func NewParser(maxComplexity int) Parser {
	if maxComplexity <= 0 {
		maxComplexity = defaultMaxComplexity
	}
	return Parser{maxComplexity: maxComplexity}
}

// DefaultParser returns a parser with the default complexity guard.
func DefaultParser() Parser {
	return Parser{maxComplexity: defaultMaxComplexity}
}

// FromJSON deserializes a JSON filter definition into an operator tree.
func (p Parser) FromJSON(payload []byte) (jsonfilter.Operator, error) {
	var root map[string]interface{}
	if err := json.Unmarshal(payload, &root); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	return p.parseRoot(root)
}

// FromYAML deserializes a YAML filter definition into an operator tree.
func (p Parser) FromYAML(payload []byte) (jsonfilter.Operator, error) {
	var root map[string]interface{}
	if err := yaml.Unmarshal(payload, &root); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	return p.parseRoot(root)
}

func (p Parser) parseRoot(node map[string]interface{}) (jsonfilter.Operator, error) {
	if node == nil {
		return nil, errors.New("filter definition cannot be empty")
	}

	if nested, ok := extractNestedMap(node, "jsonFilter"); ok {
		node = nested
	}

	op, complexity, err := p.parseOperator(node)
	if err != nil {
		return nil, err
	}
	if complexity > p.maxComplexity {
		return nil, fmt.Errorf("filter complexity %d exceeds limit %d", complexity, p.maxComplexity)
	}
	return op, nil
}

func (p Parser) parseOperator(node map[string]interface{}) (jsonfilter.Operator, int, error) {
	if len(node) == 0 {
		return nil, 0, errors.New("operator definition must contain exactly one entry")
	}
	if len(node) > 1 {
		return nil, 0, fmt.Errorf("operator definition contains multiple entries: %v", mapKeys(node))
	}

	for rawName, rawValue := range node {
		name := strings.ToLower(rawName)
		if op, count, err := p.parseComparison(name, rawValue); err == nil {
			return op, count, nil
		} else if !errors.Is(err, errUnsupportedOperator) {
			return nil, 0, err
		}

		if op, count, err := p.parseLogic(name, rawValue); err == nil {
			return op, count, nil
		} else if !errors.Is(err, errUnsupportedOperator) {
			return nil, 0, err
		}

		return nil, 0, fmt.Errorf("operator %s is not supported", rawName)
	}

	return nil, 0, errors.New("could not parse operator")
}

var errUnsupportedOperator = errors.New("unsupported operator")

func (p Parser) parseComparison(name string, value interface{}) (jsonfilter.Operator, int, error) {
	typ, err := comparison.ParseType(name)
	if err != nil {
		return nil, 0, errUnsupportedOperator
	}

	cfg, ok := normalizeMap(value)
	if !ok {
		return nil, 0, fmt.Errorf("comparison operator %s expects an object as value", name)
	}

	field, _ := cfg["field"].(string)
	if field == "" {
		return nil, 0, fmt.Errorf("comparison operator %s requires field attribute", name)
	}

	val, ok := cfg["value"]
	if !ok {
		return nil, 0, fmt.Errorf("comparison operator %s requires value attribute", name)
	}

	op, err := comparison.Instantiate(typ, field, val)
	if err != nil {
		return nil, 0, err
	}

	if v := op.Validate(); !v.Valid {
		return nil, 0, fmt.Errorf("operator %s is invalid: %s", op.Name(), v.CauseDescription)
	}

	return op, 1, nil
}

func (p Parser) parseLogic(name string, value interface{}) (jsonfilter.Operator, int, error) {
	typ, err := logic.ParseType(name)
	if err != nil {
		return nil, 0, errUnsupportedOperator
	}

	rawChildren, ok := value.([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("logic operator %s expects an array of child operators", name)
	}

	children := make([]jsonfilter.Operator, 0, len(rawChildren))
	totalComplexity := 1 // count the logic operator itself
	for idx, child := range rawChildren {
		childMap, ok := normalizeMap(child)
		if !ok {
			return nil, 0, fmt.Errorf("logic operator %s child %d must be an object", name, idx)
		}

		childOp, childComplexity, err := p.parseOperator(childMap)
		if err != nil {
			return nil, 0, err
		}
		totalComplexity += childComplexity
		if totalComplexity > p.maxComplexity {
			return nil, 0, fmt.Errorf("filter complexity %d exceeds limit %d", totalComplexity, p.maxComplexity)
		}
		children = append(children, childOp)
	}

	op, err := logic.NewOperator(typ, children)
	if err != nil {
		return nil, 0, err
	}

	if v := op.Validate(); !v.Valid {
		return nil, 0, fmt.Errorf("operator %s is invalid: %s", op.Name(), v.CauseDescription)
	}

	return op, totalComplexity, nil
}

func extractNestedMap(node map[string]interface{}, key string) (map[string]interface{}, bool) {
	if key == "" {
		return nil, false
	}
	if raw, ok := node[key]; ok {
		if nested, ok := normalizeMap(raw); ok {
			return nested, true
		}
	}
	return nil, false
}

func normalizeMap(input interface{}) (map[string]interface{}, bool) {
	switch typed := input.(type) {
	case map[string]interface{}:
		return typed, true
	case map[interface{}]interface{}:
		converted := make(map[string]interface{}, len(typed))
		for k, v := range typed {
			key, ok := k.(string)
			if !ok {
				return nil, false
			}
			converted[key] = v
		}
		return converted, true
	default:
		return nil, false
	}
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
