.PHONY: install lint test
all: lint test install demo 
install:
	go install github.com/anz-bank/sysl-catalog
lint:
	golangci-lint run ./...
test:
	go test ./...

.PHONY: demo demo-html demo-markdown
demo: demo-html demo-markdown
demo-html:
	sysl-catalog --type=html --plantuml=https://plantuml.com/plantuml -o demo/html demo/sizzle.sysl --redoc --mermaid
demo-markdown:
	sysl-catalog -o demo/markdown demo/sizzle.sysl --mermaid

.PHONY: docker docker-mermaid
docker:
	docker build -t sysl-catalog .
docker-mermaid:
	docker build -t sysl-catalog-mermaid -f mermaid.Dockerfile . 
docker-run: docker
	docker run -it -v $$(pwd)/demo/markdown:/out:rw -v $$(pwd)/demo:/usr/demo:ro anzbank/sysl-catalog demo/sizzle.sysl
docker-mermaid-run: docker-mermaid
	docker run -e SYSL_PLANTUML=localhost:8080 -v $$(pwd)/demo/html:/out:rw -v $$(pwd)/demo:/usr/demo:ro sysl-catalog-mermaid demo/sizzle.sysl --redoc --mermaid --type=html -o /out
	cp demo/mastercard.yaml demo/html/demo/mastercard.yaml
docker-compose:
	docker-compose run sysl-catalog

