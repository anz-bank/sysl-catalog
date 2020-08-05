all:
	./build.sh
install:
	go install github.com/anz-bank/sysl-catalog
lint:
	golangci-lint run ./...
.PHONY: docker build docker-mermaid
build:
	GOOS=linux GOARCH=amd64 go build -o sysl-catalog .
docker: build
	docker build -t sysl-catalog .
docker-mermaid:
	docker build -t sysl-catalog-mermaid -f mermaid.Dockerfile . 
docker-run: docker
	docker run -it -v $$(pwd)/demo/markdown:/out:rw -v $$(pwd)/demo:/usr/demo:ro anzbank/sysl-catalog demo/sizzle.sysl
docker-mermaid-run: docker-mermaid
	docker run -e SYSL_PLANTUML=localhost:8080 -v $$(pwd)/demo/html:/out:rw -v $$(pwd)/demo:/usr/demo:ro sysl-catalog-mermaid demo/sizzle.sysl --redoc --mermaid --type=html -o /out
	cp demo/mastercard.yaml demo/html/demo/mastercard.yaml
.PHONY: test
docker-compose:
	docker-compose run sysl-catalog
test:
	go test ./...
