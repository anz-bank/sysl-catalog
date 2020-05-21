all:
	rm -rf demo/html/*
	rm -rf demo/markdown/*
	rm -rf docs/*
	go run . -o demo/markdown demo/simple2.sysl
	go run . --type=html -o demo/html demo/simple2.sysl
	cp -r demo/html/* docs/
	sed -i "" "s/simple2.sysl/<a href=http:\/\/github.com\/anz-bank\/sysl-catalog>This is an example of sysl catalog deployed to github pages <\/a>/" docs/index.html
install:
	go install github.com/anz-bank/sysl-catalog

.PHONY: test
test:
	go test ./...

build:
	docker-compose build	

demo: build
	docker-compose up sysl-catalog -d
	open http://localhost:6900
