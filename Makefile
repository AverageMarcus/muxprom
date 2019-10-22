.DEFAULT_GOAL := default

.PHONY: format # Format all go code
format:
	@gofmt -s -w .

.PHONY: lint # Lint all go code
lint:
	@golint .

.PHONY: vet # Vet all code
vet:
	@go vet

.PHONY: help # Show this list of commands
help:
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1: \2/' | expand -t20

default: format lint vet
