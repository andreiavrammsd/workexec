.PHONY: all test lint coverage

GO111MODULE=on
COVER_PROFILE=cover.out

all: test lint

test:
	go test -race -cover -v ./...

lint: check-lint
	golangci-lint run

coverage:
	go test -race -v -coverprofile $(COVER_PROFILE) ./... && go tool cover -html=$(COVER_PROFILE)

prepushhook:
	echo '#!/bin/sh\n\nmake' > .git/hooks/pre-push && chmod +x .git/hooks/pre-push

check-lint:
	@[ $(shell which golangci-lint) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(shell go env GOPATH)/bin v1.50.1
