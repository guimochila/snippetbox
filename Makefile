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
