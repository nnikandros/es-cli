SHELL:=/bin/bash
ES_FIELDS_YAML:=es_fields.yaml
.PHONY: test clean build build-static



build:
	@cp $(ES_FIELDS_YAML) ./cmd/$(ES_FIELDS_YAML)
	@go build -o ./bin/es
	@mv ./bin/es $(HOME)/.local/bin


build-static:
	@cp $(ES_FIELDS_YAML) ./cmd/$(ES_FIELDS_YAML)
	@CGO_ENABLED=0 GOOS=linux go build -o ./bin/es
	@mv ./bin/es $(HOME)/.local/bin

clean:
	@rm ./bin/es


test:
	@go test ./test/

 
