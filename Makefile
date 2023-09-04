VERSION := $(shell echo "v0.0.0" | sed "s/v\([0-9.]*\)[-+]*.*/\1/g" | awk -F "." '{printf "%d.%d.%d\n",$$1,$$2,$$3+1}')

.PHONY: build
build:
	go build -o fcs-cli -trimpath -ldflags "-s -w -X main.version=${VERSION}-snapshot"

.PHONY: test
test:
	go test -cover -coverprofile cover.out

.PHONY: test-v
test-v:
	go test -v -cover -coverprofile cover.out
