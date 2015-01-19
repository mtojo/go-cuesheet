GOPATH := $(shell pwd)

lint:
	@go get github.com/golang/lint/golint
	@$(GOPATH)/bin/golint cuesheet/*.go

test:
	@cd cuesheet; go test; cd - 1>/dev/null

.PHONY: lint test
