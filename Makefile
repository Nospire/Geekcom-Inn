.PHONY: build run test clean

build:
	go build -o bin/tavrn ./cmd/tavrn

run: build
	./bin/tavrn

test:
	go test ./... -v

clean:
	rm -rf bin/ tavrn.db
