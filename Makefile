MODULE=github.com/youruser/json-filter
PKGS=./...

.PHONY: fmt lint test bench bench-logic bench-comparison tidy

fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path './vendor/*')

lint:
	go vet $(PKGS)

test:
	go test $(PKGS)

bench:
	go test ./operator/... -bench . -benchmem

bench-comparison:
	go test ./operator/comparison -bench . -benchmem


bench-logic:
	go test ./operator/logic -bench . -benchmem

tidy:
	go mod tidy
