APP=ttyslot

.PHONY: run build test tidy fmt

run:
	go run .

build:
	go build -o $(APP) .

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w *.go
