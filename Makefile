format:
	go fmt ./...
.PHONY: format

clean:
	rm -rf target/
.PHONY: target

build:
	go build -o target/depcompare cmd/cli/main.go
.PHONY: build	

test:
	go test ./...
.PHONY: test
