all: demo
.PHONY: demo

demo:
	rm -rf demo/docs/*
	rm -rf demo/docs/html/*
	go run . -o demo/markdown demo/simple2.sysl
	go run . --type=html -o demo/html demo/simple2.sysl