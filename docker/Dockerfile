# syntax=docker/dockerfile:1

FROM golang:1.18-alpine as builder

RUN apk update && apk add --no-cache g++

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY argon2id/ ./argon2id
COPY auth/ ./auth
COPY controllers/ ./controllers
COPY database/ ./database
COPY middlewares/ ./middlewares
COPY models/ ./models
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags=-static" -o go-api

FROM scratch
WORKDIR /app
ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ["/app/go-api"]
COPY --from=builder /app/go-api go-api
