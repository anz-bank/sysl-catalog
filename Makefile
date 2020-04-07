all: clean run
.PHONY: clean run demo

clean:
	rm -rf docs/*

run:
	go run . -o docs simple.sysl

demo:
	rm -rf demo/docs/*
	go run . -o demo/docs demo/simple2.sysl