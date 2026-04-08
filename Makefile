.PHONY: tidy fmt test tests

PKGS_WITH_TESTS := $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

fmt:
	gofumpt -extra -l -w .

tidy:
	go mod tidy \

test:
	go test $(PKGS_WITH_TESTS) -cover
