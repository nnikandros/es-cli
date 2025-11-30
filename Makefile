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

completion: build
	@es completion bash > es_bash.sh
	@source es_bash.sh

.PHONY: restore
restore:
	@complete -r es
	@rm es_bash.sh

