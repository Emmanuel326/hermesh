APP=hermesh
VERSION=0.1.0
GOEXPERIMENT=greenteagc

.PHONY: run build clean build-all

run:
	go run ./...

build:
	GOEXPERIMENT=$(GOEXPERIMENT) go build -ldflags="-s -w -X main.Version=$(VERSION)" -o bin/$(APP) ./...

build-all:
	mkdir -p bin
	GOEXPERIMENT=$(GOEXPERIMENT) GOOS=linux   GOARCH=amd64  go build -ldflags="-s -w" -o bin/$(APP)-linux-amd64 ./...
	GOEXPERIMENT=$(GOEXPERIMENT) GOOS=linux   GOARCH=arm64  go build -ldflags="-s -w" -o bin/$(APP)-linux-arm64 ./...
	GOEXPERIMENT=$(GOEXPERIMENT) GOOS=darwin  GOARCH=amd64  go build -ldflags="-s -w" -o bin/$(APP)-darwin-amd64 ./...
	GOEXPERIMENT=$(GOEXPERIMENT) GOOS=darwin  GOARCH=arm64  go build -ldflags="-s -w" -o bin/$(APP)-darwin-arm64 ./...
	GOEXPERIMENT=$(GOEXPERIMENT) GOOS=windows GOARCH=amd64  go build -ldflags="-s -w" -o bin/$(APP)-windows-amd64.exe ./...

clean:
	rm -rf bin/

tidy:
	go mod tidy

lint:
	go vet ./...
