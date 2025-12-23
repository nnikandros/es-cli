SHELL:=/bin/bash
ES_FIELDS_YAML:=es_fields.yaml
.PHONY: test clean

hello:
	@if [-z $(ES_FIELDS_YAML)]; then\

	echo $(ES_FIELDS_YAML)
	cp $(ES_FIELDS_YAML) ./cmd/$(ES_FIELDS_YAML)
	stat ./cmd/$(ES_FIELDS_YAML)
	rm ./cmd/$(ES_FIELDS_YAML)

build:
	@cp $(ES_FIELDS_YAML) ./cmd/$(ES_FIELDS_YAML)
	@go build -o ./bin/es
	@mv ./bin/es $(HOME)/.local/bin
	@rm ./cmd/$(ES_FIELDS_YAML)


build-static:
	@CGO_ENABLED=0 GOOS=linux go build -o ./bin/es

clean:
	@rm ./bin/es


test:
	@go test ./test/

 
 check-var:
	@if [ -z "$(ES_FIELDS_YAML)" ]; then \
		echo "Error: MY_VAR is not set"; \
		exit 1; \
	else \
		echo "MY_VAR is $(ES_FIELDS_YAML)"; \
	fi