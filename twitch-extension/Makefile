.PHONY: openapi deps

# Installs all dependencies for the project.
deps:
	cd ebs && go mod tidy && \
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0

# Auto-generates the OpenAPI client and server code from the yaml file.
openapi:
	mkdir -p ./ebs/openapi
	oapi-codegen -generate gin,models -package openapi ./openapi/openapi.yaml > ./ebs/openapi/api.gen.go
