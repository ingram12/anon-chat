.PHONY: generate-api lint test build build-frontend build-backend

generate-api:
	mkdir -p internal/api
	oapi-codegen --generate types -o internal/api/models.gen.go -package api api/openapi/openapi.json
	oapi-codegen --generate server -o internal/api/server.gen.go -package api api/openapi/openapi.json
	oapi-codegen --generate spec -o internal/api/spec.gen.go -package api api/openapi/openapi.json

lint:
	golangci-lint run ./...

test:
	go test -v -count=1 ./...

build-frontend:
	npm install --prefix frontend
	npm run build --prefix frontend
	rm -rf ./build/frontend
	mkdir -p ./build/frontend
	cp -r frontend/dist/. ./build/frontend

build-backend:
	go mod download
	go build -o build/server ./cmd/server

build: build-frontend build-backend
