.PHONY: install lint tidy test test-integration coverage
all: lint test install demo 
install:
	go install github.com/anz-bank/sysl-catalog
lint:
	golangci-lint run ./...
tidy:
	go mod tidy
	gofmt -s -w .
	goimports -w .
coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./... && go tool cover -func=coverage.txt
test:
	go test ./...
test-integration: # Run integration tests (against Github API. Requires SYSL_GITHUB_TOKEN to be set)
	go test -tags=integration ./...

.PHONY: demo demo-html demo-markdown
demo: demo-html demo-markdown
demo-html:
	sysl-catalog --type=html -o demo/html demo/demo.sysl --redoc --mermaid
demo-markdown:
	sysl-catalog -o demo/markdown demo/demo.sysl --mermaid
demo-server:
	docker run \
	-p 6900:6900 \
	-e SYSL_GITHUB_TOKEN=$(SYSL_GITHUB_TOKEN) \
	-e SYSL_PLANTUML=localhost:8080 \
	-e SYSL_MODULES=github \
	-v $$(pwd)/demo/html:/out:rw \
	-v $$(pwd)/demo:/usr/demo:ro \
	sysl-catalog-mermaid \
		--serve \
		-v \
		--redoc \
		--mermaid \
		--type=html \
		-o /out \
		demo/demo.sysl
.PHONY: docker docker-mermaid
docker:
	docker build -t sysl-catalog .
docker-mermaid:
	docker build -t sysl-catalog-mermaid -f mermaid.Dockerfile . 
docker-run: docker
	docker run -it -e SYSL_PLANTUML=localhost:8080 -e SYSL_MODULES=github -v $$(pwd)/demo/markdown:/out:rw -v $$(pwd)/demo:/usr/demo:ro anzbank/sysl-catalog demo/demo.sysl
docker-mermaid-run: docker-mermaid
	docker run -e SYSL_PLANTUML=localhost:8080 -e SYSL_MODULES=github -v $$(pwd)/demo/html:/out:rw -v $$(pwd)/demo:/usr/demo:ro sysl-catalog-mermaid demo/demo.sysl --redoc --mermaid --type=html -o /out
docker-compose:
	docker-compose run sysl-catalog
