VERSION := $(shell git describe --tags --abbrev=0 | awk -F "." '{sub("v","", $$1); printf "%s.%s.%s\n",$$1,$$2,$$3+1}')

BINARY=fcs-cli
GO_FILES=main.go
GOCOVERDIR=coverdir

$(BINARY): $(GO_FILES)
	go build -o $@ -trimpath -ldflags "-s -w -X main.version=${VERSION}-snapshot"

test/$(BINARY): $(GO_FILES)
	go build -o $@ -cover -trimpath -ldflags "-s -w -X main.version=0.0.0-test"

.PHONY: test
test:
	go test -cover -coverprofile cover-ut.out
	@go tool cover -html=cover-ut.out -o cover-ut.html

.PHONY: test-v
test-v:
	go test -v -cover -coverprofile cover-ut.out
	@go tool cover -html=cover-ut.out -o cover-ut.html

.PHONY: integration-test
integration-test: test/$(BINARY)
	@mkdir -p coverdir
	GOCOVERDIR=coverdir go test -v ./test
	@go tool covdata percent -i=$(GOCOVERDIR)
	@go tool covdata textfmt -i=$(GOCOVERDIR) -o cover-it.out
	@go tool cover -html=cover-it.out -o cover-it.html

.PHONY: clean
clean:
	rm -rf $(BINARY) test/$(BINARY) dist cover*
