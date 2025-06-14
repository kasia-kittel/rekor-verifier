test:
	go test ./...

build:
	go build -o bin/main main.go

deps:
	go get

run:
	go run main.go