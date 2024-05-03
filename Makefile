BINARY_NAME=main
.DEFAULT_GOAL := run

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} cmd/web/main.go

run: build
	./${BINARY_NAME}

dep:
	go mod download

clean:
	go clean
	rm ${BINARY_NAME}
	rm tmp/

lint:
	golangci-lint run --enable-all

test:
	go test -v -short ./...

test-all:
	go test -v ./...

test-coverage:
	rm -rf tmp/coverage
	mkdir tmp/coverage
	go test -covermode=count -coverprofile=tmp/coverage/profile.out ./...
	go tool cover -html=tmp/coverage/profile.out
