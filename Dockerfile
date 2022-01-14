# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY app/ .

RUN go mod download
RUN go get github.com/google/uuid
ENV CGO_ENABLED 0
RUN go build -o /go_api

ENTRYPOINT  go test ./tests && /go_api && /bin/bash