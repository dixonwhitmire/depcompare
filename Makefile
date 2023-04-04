format:
	go fmt ./...
.PHONY: format

clean:
	rm -rf target/
.PHONY: target

build:
	go build -o target/app cmd/cli/main.go
.PHONY: build	

test:
	go test ./...
.PHONY: test

run:
	go run cmd/cli/main.go
.PHONY: run