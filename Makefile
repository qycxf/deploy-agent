MODULE_NAME=github.com/qycxf/deploy-agent
BINARY=deploy-agent

.PHONY: all gen build run tidy fmt test docker-build bundle validate gen-oapi gen-openapi

all: build

tidy:
	go mod tidy

fmt:
	go fmt ./...

bundle:
	# prefer local install in node_modules/.bin, otherwise fallback to npx
	@if [ -x ./node_modules/.bin/swagger-cli ]; then \
		./node_modules/.bin/swagger-cli bundle open-api/openapi.yaml -o open-api/bundled.yaml -t yaml; \
	else \
		npx --yes @apidevtools/swagger-cli bundle open-api/openapi.yaml -o open-api/bundled.yaml -t yaml; \
	fi

validate:
	# prefer local binaries, fallback to npx
	@if [ -x ./node_modules/.bin/swagger-cli ]; then \
		./node_modules/.bin/swagger-cli validate open-api/bundled.yaml; \
	else \
		npx --yes @apidevtools/swagger-cli validate open-api/bundled.yaml; \
	fi
	@if [ -x ./node_modules/.bin/spectral ]; then \
		./node_modules/.bin/spectral lint open-api/bundled.yaml || true; \
	else \
		npx --yes @stoplight/spectral lint open-api/bundled.yaml || true; \
	fi

tools-install:
	# install local cli tools to node_modules/.bin (no sudo, reusable)
	npm install --no-save @apidevtools/swagger-cli @stoplight/spectral
	# install Go-based generators (oapi-codegen v2)
	@command -v go >/dev/null 2>&1 && go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest || true

tools-clean:
	rm -rf node_modules package-lock.json

.PHONY: generate

generate:
	@echo "🚀 generate model (Model)..."
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config internal/model/cfg.yaml open-api/bundled.yaml
	@echo "🚀 generate server (API)..."
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config internal/api/cfg.yaml open-api/bundled.yaml
	@echo "✅ All code generation completed successfully!"

gen-openapi:
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/open-api/bundled.yaml -g go-gin-server -o /local/gen/openapi-generator --additional-properties=packageName=api

build:
	# build the root CLI binary (optional)
	go build -v -o bin/$(BINARY) .

run:
	# run the server subcommand via the root CLI
	go run . server

test:
	go test ./...

docker-build:
	docker build -t $(BINARY):latest .
