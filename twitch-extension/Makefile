.PHONY: openapi deps ebs

# Installs all dependencies for the project.
deps:
	cd ebs && go mod tidy && \
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0

# Auto-generates the OpenAPI client and server code from the yaml file.
openapi:
	mkdir -p ./ebs/server/openapi
	oapi-codegen -generate gin,models -package openapi ./openapi/openapi.yaml > ./ebs/server/openapi/api.gen.go

ebs:
	@trap 'kill 0' INT TERM; \
	cd frontend && npm start & \
	cd ebs && go run . -disableAuthentication