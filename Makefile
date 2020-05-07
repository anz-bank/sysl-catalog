all:
	rm -rf demo/html/*
	rm -rf demo/markdown/*
	rm -rf docs/*
	go run . -o demo/markdown demo/simple2.sysl
	go run . --type=html -o demo/html demo/simple2.sysl
	cp -r demo/html/* docs/
	printf '%s\n%s\n' "<a href=\"http://github.com/anz-bank/sysl-catalog\">This is an example of sysl catalog deployed to github pages </a>" "$$(cat docs/index.html)" > docs/index.html
install:
	go install github.com/anz-bank/sysl-catalog
