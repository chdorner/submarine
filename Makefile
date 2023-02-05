.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: build
build: # Build submarine
	go build ./...

.PHONY: test
test: # Test submarine
	go test -cover ./...

.PHONY: coverage
coverage: # Run tests and open coverage report
	mkdir -p tmp
	go test -coverprofile=tmp/coverage.out ./...
	go tool cover -html=tmp/coverage.out

.PHONY: install-tools
install-tools: # Install development tools
	brew install goreleaser/tap/goreleaser
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1

.PHONY: lint
lint: # Lint source code
	golangci-lint run
