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
		./node_modules/.bin/swagger-cli bundle api/openapi.yaml -o api/bundled.yaml -t yaml; \
	else \
		npx --yes @apidevtools/swagger-cli bundle api/openapi.yaml -o api/bundled.yaml -t yaml; \
	fi

validate:
	# prefer local binaries, fallback to npx
	@if [ -x ./node_modules/.bin/swagger-cli ]; then \
		./node_modules/.bin/swagger-cli validate api/bundled.yaml; \
	else \
		npx --yes @apidevtools/swagger-cli validate api/bundled.yaml; \
	fi
	@if [ -x ./node_modules/.bin/spectral ]; then \
		./node_modules/.bin/spectral lint api/bundled.yaml || true; \
	else \
		npx --yes @stoplight/spectral lint api/bundled.yaml || true; \
	fi

tools-install:
	# install local cli tools to node_modules/.bin (no sudo, reusable)
	npm install --no-save @apidevtools/swagger-cli @stoplight/spectral
	# install Go-based generators (oapi-codegen v2)
	@command -v go >/dev/null 2>&1 && go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v2 || true

tools-clean:
	rm -rf node_modules package-lock.json

gen-oapi:
	oapi-codegen --config=oapi-codegen.yaml open-api/bundled.yaml

gen-openapi:
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/api/bundled.yaml -g go-gin-server -o /local/gen/openapi-generator --additional-properties=packageName=api

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
