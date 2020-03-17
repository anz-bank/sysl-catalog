all: clean run
.PHONY: clean run

clean:
	rm -rf docs/*

run:
	go run . -o docs simple.sysl