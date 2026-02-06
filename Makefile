.PHONY: test test-short test-cover lint vet build examples clean

test:
	go test ./... -race -count=1

test-short:
	go test ./... -short

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

vet:
	go vet ./...

build:
	go build ./...

examples:
	go build ./examples/...

clean:
	rm -f coverage.out coverage.html
