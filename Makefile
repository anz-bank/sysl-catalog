all:
	git rm -rf demo/html/* || true 
	git rm -rf demo/markdown/* || true 
	git rm -rf docs/* || true 
	go run . --plantuml=https://plantuml.com/plantuml --embed -o demo/markdown demo/simple2.sysl
	go run . --plantuml=https://plantuml.com/plantuml --type=html --embed -o demo/html demo/simple2.sysl
	mkdir -p docs
	cp -r demo/html/* docs/
	sed -i "" "s/simple2.sysl/<a href=http:\/\/github.com\/anz-bank\/sysl-catalog>This is an example of sysl catalog deployed to github pages <\/a>/" docs/index.html
install:
	go install github.com/anz-bank/sysl-catalog

.PHONY: test
test:
	go test ./...
