NAME     := ghrls
VERSION  := v0.1.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"github.com/dtan4/ghrls/version.Version=$(VERSION)\" -X \"github.com/dtan4/ghrls/version.Revision=$(REVISION)\" -extldflags \"-static\""
NOVENDOR := $(shell go list ./... | grep -v vendor)

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-test
ci-test:
	go test -coverpkg=./... -coverprofile=coverage.txt -v ./...

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: cross-build
cross-build:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: test
test:
	go test -cover -v $(NOVENDOR)
