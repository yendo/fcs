VERSION := $(shell cat VERSION )

.PHONY: build
build:
	go build -o fcs-cli -trimpath -ldflags "-s -w -X main.version=${VERSION}"

.PHONY: test
test:
	go test -cover -coverprofile cover.out

.PHONY: test-v
test-v:
	go test -v -cover -coverprofile cover.out
