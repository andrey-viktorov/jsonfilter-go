# JSON-Filter Go Implementation

## Project Overview

This is a Go port of [telekom/JSON-Filter](https://github.com/telekom/JSON-Filter.git), a lightweight library for evaluating JSON payloads against filter operators. This library will be integrated into a high-performance mocking REST server.

**Performance Requirements:**
- Use `fasthttp` for zero-allocation HTTP handling
- Work with `[]byte` instead of strings where possible
- Use `gjson` library for JSON path field selection (equivalent to Java's JsonPath)
- Optimize for minimal allocations and maximum throughput

## Architecture

### Core Components

1. **Operator Interface** (`operator.go`)
   - Base interface with `Evaluate(json string) EvaluationResult` and `Validate() ValidationResult` methods
   - Two operator categories: Comparison and Logic

2. **Comparison Operators** (in `operator/comparison/`)
   - Equal (`eq`), NotEqual (`ne`), Regex (`rx`)
   - LessThan (`lt`), LessEqual (`le`), GreaterThan (`gt`), GreaterEqual (`ge`)
   - In (`in`), NotIn (`nin`), Contains (`ct`), NotContains (`nct`)
   - Each operator has: `field` (JSON path), `value` (expected value), `operator` (enum)
   - Use `gjson.Get()` and `gjson.GetMany()` for path evaluation

3. **Logic Operators** (in `operator/logic/`)
   - `and`: Valid if ALL child operators are valid
   - `or`: Valid if AT LEAST ONE child operator is valid
   - Contains list of child operators

4. **EvaluationResult** (`evaluation_result.go`)
   - Fields: `Match` (bool), `OperatorName` (string), `CauseDescription` (string), `ChildOperators` ([]EvaluationResult)
   - Factory methods: `Valid()`, `WithError()`, `FromResultList()` (for logic operators)

5. **Serialization/Deserialization** (`serde/`)
   - Support JSON and YAML operator definitions
   - Implement custom unmarshaling for operator chains
   - Validate operator complexity (default max: 42) to prevent abuse

### Data Flow

```
YAML/JSON Filter Definition → Unmarshal → Operator Tree → Evaluate(json) → EvaluationResult
```

## Development Patterns

### Using gjson for JSON Path Evaluation

```go
// Single value extraction (for eq, ne, lt, etc.)
result := gjson.Get(json, "$.processing.state")
if !result.Exists() {
    return EvaluationResult.WithError("path not found")
}

// Multiple values (for contains, not-contains)
results := gjson.Get(json, "$..items").Array()
```

### Operator Implementation Template

```go
type EqualOperator struct {
    operator    ComparisonOperatorEnum
    jsonPath    string
    expectedValue interface{}
}

func (o *EqualOperator) Evaluate(json string) EvaluationResult {
    actualValue := gjson.Get(json, o.jsonPath)
    if actualValue.Value() == o.expectedValue {
        return EvaluationResult.Valid(o)
    }
    return EvaluationResult.WithError(o, "values don't match")
}
```

### Performance Considerations

- **Avoid string concatenation** in hot paths - use `strings.Builder`
- **Reuse gjson results** - don't parse the same path multiple times
- **Lazy evaluation** for `or` operators - stop on first match
- **Pre-compile regex patterns** in RegexOperator validation phase

### Error Handling

- Return detailed error messages in `EvaluationResult.CauseDescription`
- Include operator context and actual vs expected values
- For invalid JSON paths, catch in `Validate()` before `Evaluate()`

## Testing Strategy

Reference the Java test files at `src/test/java/de/telekom/jsonfilter/`:
- `OperatorTest.java` - End-to-end filter evaluation tests
- `OperatorDeserializerTest.java` - YAML/JSON parsing tests
- Test with `serdeTest.yaml`, `validPayload.json`, `nctBlacklist.yaml`

Test categories:
1. Individual operator correctness (each comparison type)
2. Logic operator combinations (nested and/or)
3. Complex filter chains (validate complexity limits)
4. Edge cases (missing fields, type mismatches, array handling)

## Common Tasks

### Adding a New Comparison Operator

1. Add enum to `ComparisonOperatorEnum`
2. Create operator file in `operator/comparison/`
3. Implement `Evaluate()` and `Validate()` methods
4. Add factory case in `ComparisonOperator.Instantiate()`
5. Update deserializer to recognize operator key
6. Add test cases with YAML examples

### Deserializing from YAML

```go
type FilterConfig struct {
    JsonFilter Operator `yaml:"jsonFilter"`
}

// Implement custom UnmarshalYAML for Operator interface
```

### Integration with Mocking Server

This library provides the filtering core. The mocking server should:
- Accept YAML config with filters and mock responses
- Use `fasthttp` for request handling
- Extract request body as `[]byte`
- Convert to string only when calling `operator.Evaluate()`
- Match filters and return configured mock response

## Key Files from Reference Implementation

- `src/main/java/de/telekom/jsonfilter/operator/Operator.java` - Base interface
- `src/main/java/de/telekom/jsonfilter/operator/comparison/ComparisonOperator.java` - Comparison base with gjson equivalent
- `src/main/java/de/telekom/jsonfilter/operator/logic/LogicOperator.java` - Logic base
- `src/main/java/de/telekom/jsonfilter/serde/OperatorDeserializer.java` - Complex recursive deserialization logic
- `src/main/java/de/telekom/jsonfilter/operator/EvaluationResult.java` - Result structure

## Go-Specific Conventions

- Use Go modules (`go.mod`)
- Package structure: `github.com/andrey-viktorov/jsonfilter-go`
- Dependencies: `github.com/tidwall/gjson`, `github.com/valyala/fasthttp`, `gopkg.in/yaml.v3`
- Exported types start with capital letters
- Use interfaces for polymorphism (Operator, OperatorEnum)
- Prefer value receivers unless mutating state
- Use table-driven tests with `testing` package
