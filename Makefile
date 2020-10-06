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
	sysl-catalog run --type=html -o demo/html demo/demo.sysl
demo-markdown:
	sysl-catalog run -o demo/markdown demo/demo.sysl
demo-server:
	docker run \
	-p 6900:6900 \
	-e SYSL_TOKENS=$(SYSL_TOKENS) \
	-e SYSL_PLANTUML=localhost:8080 \
	-e SYSL_MODULES=on \
	-v $$(pwd)/demo/html:/out:rw \
	-v $$(pwd)/demo:/usr/demo:ro \
	sysl-catalog run \
		--serve \
		-v \
		--type=html \
		-o /out \
		demo/demo.sysl
.PHONY: docker
docker:
	docker build -t sysl-catalog .
docker-run: docker
	docker run -it -e SYSL_PLANTUML=localhost:8080 -e SYSL_MODULES=on -v $$(pwd)/demo/markdown:/out:rw -v $$(pwd)/demo:/usr/demo:ro anzbank/sysl-catalog run demo/demo.sysl
docker-compose:
	docker-compose run sysl-catalog
