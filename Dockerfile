FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ ./
RUN go build -o app .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM ubuntu:latest
WORKDIR /usr/local/bin
COPY --from=builder /build/app .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /build/migrations /migrations
COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh