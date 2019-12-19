.PHONY: all test lint coverage

GO111MODULE=on
COVER_PROFILE=cover.out

all: test lint

test:
	go test -race -cover -v ./...

lint: check-lint
	golint -set_exit_status ./...
	golangci-lint run

coverage:
	go test -race -v -coverprofile $(COVER_PROFILE) ./... && go tool cover -html=$(COVER_PROFILE)

prepushhook:
	echo '#!/bin/sh\n\nmake' > .git/hooks/pre-push && chmod +x .git/hooks/pre-push

check-lint:
	@[ $(shell which golint) ] || (GO111MODULE=off && go get -u golang.org/x/lint/golint)
	@[ $(shell which golangci-lint) ] || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh \
		| sh -s -- -b $(shell go env GOPATH)/bin v1.21.0
