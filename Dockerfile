FROM golang:1.22.3-alpine as build-stage

WORKDIR /app

COPY go.mod .
RUN go mod download

# Copy source code
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /snippetbox ./cmd/web

FROM alpine:latest as build-release

WORKDIR /app

COPY --from=build-stage /snippetbox /snippetbox

ENTRYPOINT ["/snippetbox"]
