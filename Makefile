all:
	git rm -rf demo/html/* || true 
	git rm -rf demo/markdown/* || true 
	git rm -rf docs/* || true 
	go run . -o demo/markdown demo/simple2.sysl
	go run . --type=html -o demo/html demo/simple2.sysl
	cp -r demo/html/* docs/
	sed -i "" "s/simple2.sysl/<a href=http:\/\/github.com\/anz-bank\/sysl-catalog>This is an example of sysl catalog deployed to github pages <\/a>/" docs/index.html
install:
	go install github.com/anz-bank/sysl-catalog

.PHONY: test
test:
	go test ./...
