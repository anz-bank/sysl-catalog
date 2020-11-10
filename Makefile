.PHONY: install lint tidy test test-integration coverage
all: lint test install demo 
install:
	go install github.com/anz-bank/sysl-catalog
lint:
	golangci-lint run ./...
tidy:
	goimports -w .
	gofmt -s -w .
	go mod tidy
coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./... && go tool cover -func=coverage.txt
test:
	go test ./...
test-integration: # Run integration tests (against Github API. Requires SYSL_GITHUB_TOKEN to be set)
	go test -tags=integration ./...

.PHONY: demo demo-html demo-markdown
demo: demo-html demo-markdown
demo-html:
	sysl-catalog run --type=html -o demo/html demo/demo.sysl
demo-markdown:
	sysl-catalog run -o demo/markdown demo/demo.sysl
demo-docker: docker
	docker run \
	-e SYSL_TOKENS=$(SYSL_TOKENS) \
	-e SYSL_MODULES=on \
	-v $$(pwd):/usr/ \
	anzbank/sysl-catalog:latest run \
		-o demo/html \
		-v \
		--type=html \
		demo/demo.sysl
demo-server: docker
	docker run \
	-p 6900:6900 \
	-e SYSL_TOKENS=$(SYSL_TOKENS) \
	-e SYSL_PLANTUML=localhost:8080 \
	-e SYSL_MODULES=on \
	-v $$(pwd)/demo:/usr/demo:ro \
	anzbank/sysl-catalog:latest run \
		--serve \
		-v \
		--type=html \
		demo/demo.sysl
.PHONY: docker
docker:
	docker build -t anzbank/sysl-catalog:latest .
