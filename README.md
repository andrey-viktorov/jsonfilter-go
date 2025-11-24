JSON-Filter (Go)
=================

[![pkg.go.dev](https://pkg.go.dev/badge/github.com/andrey-viktorov/jsonfilter-go.svg)](https://pkg.go.dev/github.com/andrey-viktorov/jsonfilter-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

High-performance Go port of [telekom/JSON-Filter](https://github.com/telekom/JSON-Filter). The library evaluates JSON payloads against declarative filter definitions (JSON or YAML) and returns detailed match results suitable for ultra-fast mocking servers built on `fasthttp`.

Features
--------

- **Zero-allocation hot paths** – comparison and logic operators operate directly on `[]byte` payloads via `gjson` without extra copies.
- **Rich operator set** – equality, regex, and logic (`and`, `or`) operators implemented with the same semantics as the reference project. Additional comparison operators can be added via the shared factory.
- **Serde with complexity guards** – load filters from JSON or YAML, enforce a configurable max tree complexity (default 42) to prevent abuse.
- **Detailed evaluation and validation results** – every operator can validate itself before execution and produce structured match reports.
- **Benchmarked** – reproducible Go benchmarks document latency and allocation characteristics.

Project Layout
--------------

```
.
├── operator
│   ├── comparison   # eq/rx operators, factories, tests, benchmarks
│   └── logic        # and/or operator implementation, tests, benchmarks
├── serde            # Parser for JSON/YAML filter definitions + tests
├── evaluation_result.go / validation_result.go
├── operator.go      # Operator interface shared across packages
└── Makefile         # Formatting, linting, testing, benchmarking helpers
```

Getting Started
---------------

```bash
git clone https://github.com/andrey-viktorov/jsonfilter-go
cd jsonfilter-go
go test ./...
```

Usage Example
-------------

```go
package main

import (
	"fmt"

	jsonfilter "github.com/andrey-viktorov/jsonfilter-go"
	"github.com/andrey-viktorov/jsonfilter-go/serde"
)

func main() {
	parser := serde.DefaultParser()
	filterYAML := []byte(`
jsonFilter:
  and:
	- eq:
		field: $.processing.state
		value: done
	- rx:
		field: $.payload.id
		value: ^[A-Z]{3}-[0-9]{4}$
`)

	op, err := parser.FromYAML(filterYAML)
	if err != nil {
		panic(err)
	}

	body := []byte(`{"processing":{"state":"done"},"payload":{"id":"ABC-1234"}}`)
	result := op.Evaluate(body)
	fmt.Printf("Match=%v, Cause=%s\n", result.Match, result.CauseDescription)
}
```

Serde Format
------------

Filters use a single root operator. Each comparison operator requires `field` (JSON path understood by `gjson`) and `value`.

```yaml
jsonFilter:
  or:
    - eq:
      field: $.metadata.status
      value: ready
    - rx:
      field: $.metadata.traceId
      value: "^trace-[0-9]+$"
```

Complexity Guard
----------------

```go
parser := serde.NewParser(10) // disallow filters with complexity > 10
```

Benchmarks (Apple M4, Go 1.21)
------------------------------

| Benchmark | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| Equal match | 27.06 | 0 | 0 |
| Equal miss  | 27.12 | 0 | 0 |
| Regex match | 82.27 | 0 | 0 |
| And eval    | 87.43 | 0 | 0 |
| Or eval     | 34.56 | 0 | 0 |

Testing
-------

```bash
make test        # go test ./...
make bench       # go test ./operator/... -bench . -benchmem
```

Contributing
------------

1. Fork and branch.
2. Add/extend operators or serde logic.
3. Run `make fmt lint test`.
4. Update docs/benchmarks when behavior changes.

License
-------

MIT. See [LICENSE](LICENSE) for details.