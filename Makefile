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
	docker build . -t anz-bank/sysl-catalog

demo: build
	docker run -p 6900:6900 --entrypoint=sh -v $(pwd)/demo:/demo anz-bank/sysl-catalog -c '/usr/src/sysl-catalog /demo/simple2.sysl --serve' -d