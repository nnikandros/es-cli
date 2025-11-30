SHELL:=/bin/bash

build:
	@go build -o ./bin/es
	@mv ./bin/es ../bin


build-static:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/es
	@mv ./bin/es ../bin


.PHONY: test
test:
	@go test ./test/

