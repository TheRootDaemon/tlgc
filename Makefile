# Project Packages
PKGS = $(shell go list ./...)
PKGS_WITH_TESTS := $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")

#
# Build targets
#

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build \
		-trimpath \
		-ldflags="-s -w -X github.com/TheRootDaemon/tlgc/version.Version=$(VERSION)" \
		-o bin/tldr \
		./main.go

.PHONY: install
install:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go install \
		-trimpath \
		-ldflags="-s -w -X github.com/TheRootDaemon/tlgc/version.Version=$(VERSION)" \
		.

#
# Development targets
#

.PHONY: run
run: 
	go run ./main.go

.PHONY: tidy
tidy:
	go mod tidy

#
# Quality targets
#

.PHONY: check
check: fmt lint sec test vet

.PHONY: cover
cover:
	go test -race \
	$(PKGS_WITH_TESTS) \
	-coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: cover-html
cover-html:
	go test -race \
	$(PKGS_WITH_TESTS) \
	-coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: fmt
fmt:
	gofumpt -extra -l -w .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: sec
sec:
	gosec \
	-exclude-generated \
	./...

.PHONY: test
test:
	go test -race \
	$(PKGS_WITH_TESTS) -cover

.PHONY: vet
vet:
	go vet $(PKGS)

#
# Maintanence targets
#

.PHONY: clean
clean:
	rm -rf bin coverage.out
