SHELL:=/bin/bash

build:
	@go build -o ./bin/es
	@mv ./bin/es $(HOME)/.local/bin


build-static:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/es
	@mv ./bin/es $(HOME)/.local/bin


.PHONY: test
test:
	@go test ./test/

