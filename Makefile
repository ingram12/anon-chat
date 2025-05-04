.PHONY: generate-api lint test build build-frontend build-backend

generate-api:
	mkdir -p internal/api
	oapi-codegen --generate types -o internal/api/models.gen.go -package api api/openapi/openapi.json
	oapi-codegen --generate server -o internal/api/server.gen.go -package api api/openapi/openapi.json
	oapi-codegen --generate spec -o internal/api/spec.gen.go -package api api/openapi/openapi.json

lint:
	golangci-lint run ./...

test:
	go test -v ./...

build-frontend:
	cd frontend && npm install && npm run build

build-backend:
	go build -o bin/server ./cmd/server

build: build-frontend build-backend
