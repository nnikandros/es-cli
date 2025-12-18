SHELL:=/bin/bash
.PHONY: test clean

build:
	@go build -o ./bin/es
	@mv ./bin/es $(HOME)/.local/bin


build-static:
	@CGO_ENABLED=0 GOOS=linux go build -o ./bin/es

clean:
	@rm ./bin/es


test:
	@go test ./test/

