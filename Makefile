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
	docker run -v $$(pwd)/demo/markdown:/out:rw -v $$(pwd)/demo:/usr/demo:ro anzbank/sysl-catalog demo/simple2.sysl
.PHONY: test
docker-compose:
	docker-compose run sysl-catalog
test:
	go test ./...
