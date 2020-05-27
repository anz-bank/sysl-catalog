all:
	./build.sh

install:
	go install github.com/anz-bank/sysl-catalog

.PHONY: docker

docker:
	docker build -t sysl-catalog .

docs: docker
	(PLANTUML_IMAGE=plantuml/plantuml-server:jetty docker-compose run docs)

.PHONY: test
test:
	go test ./...
