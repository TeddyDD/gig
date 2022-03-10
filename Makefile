GO_FILES = $(shell find . -type f -iname '*.go')

coverage.txt: $(GO_FILES)
	go test -coverprofile=coverage.txt -covermode=atomic -timeout 5s .

.PHONY: lint
lint:
	golangci-lint run

test: coverage.txt
	go test -race .
	$(MAKE) lint

bench:
	go test -run=nothing -bench=.

html: coverage.txt
	go tool cover -html=coverage.txt
