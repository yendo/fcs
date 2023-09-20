VERSION := $(shell git describe --tags --abbrev=0 | awk -F "." '{sub("v","", $$1); printf "%s.%s.%s\n",$$1,$$2,$$3+1}')

BINARY=fcs-cli
GO_FILES=main.go go.mod go.sum
GOCOVERDIR=coverdir

$(BINARY): $(GO_FILES)
	go build -o $@ -trimpath -ldflags "-s -w -X main.version=${VERSION}-snapshot"

test/$(BINARY): $(GO_FILES)
	go build -o $@ -cover -trimpath -ldflags "-s -w -X main.version=0.0.0-test"

.PHONY: test
test: unit-test integration-test

.PHONY: test-v
test-v: unit-test-v integration-test-v

.PHONY: unit-test
unit-test:
	go test -cover -coverprofile cover-ut.out
	@go tool cover -html=cover-ut.out -o cover-ut.html

.PHONY: unit-test-v
unit-test-v:
	go test -v -cover -coverprofile cover-ut.out
	@go tool cover -html=cover-ut.out -o cover-ut.html

.PHONY: integration-test
integration-test: test/$(BINARY)
	@mkdir -p coverdir
	GOCOVERDIR=coverdir go test ./test
	@go tool covdata percent -i=$(GOCOVERDIR)
	@go tool covdata textfmt -i=$(GOCOVERDIR) -o cover-it.out
	@go tool cover -html=cover-it.out -o cover-it.html

.PHONY: integration-test-v
integration-test-v: test/$(BINARY)
	@mkdir -p coverdir
	GOCOVERDIR=coverdir go test -v ./test
	@go tool covdata percent -i=$(GOCOVERDIR)
	@go tool covdata textfmt -i=$(GOCOVERDIR) -o cover-it.out
	@go tool cover -html=cover-it.out -o cover-it.html

.PHONY: super-linter
super-linter: clean
	docker run -e RUN_LOCAL=true -e USE_FIND_ALGORITHM=true -e VALIDATE_GO=false \
		-v ${PWD}:/tmp/lint/ ghcr.io/super-linter/super-linter:slim-v5

.PHONY: clean
clean:
	rm -rf $(BINARY) test/$(BINARY) dist cover*
