.PHONY: generate-api

generate-api:
	mkdir -p internal/api
	oapi-codegen --generate types -o internal/api/models.gen.go -package api api/openapi/openapi.json
	oapi-codegen --generate server -o internal/api/server.gen.go -package api api/openapi/openapi.json
	oapi-codegen --generate spec -o internal/api/spec.gen.go -package api api/openapi/openapi.json
