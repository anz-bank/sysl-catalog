all:
	./build.sh

install:
	go install github.com/anz-bank/sysl-catalog

.PHONY: docker build
build:
	GOOS=linux GOARCH=amd64 go build -o sysl-catalog .
docker: build
	docker build -t sysl-catalog .
docker-run: docker
	docker run -v $$(pwd)/this:/out:rw -v $$(pwd)/demo:/usr/src/demo:ro sysl-catalog
docs: docker
	(PLANTUML_IMAGE=plantuml/plantuml-server:jetty docker-compose run docs)

.PHONY: test
test:
	go test ./...
