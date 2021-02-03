FROM golang:1.15-alpine

LABEL mantainer="Auras Popescu popescuauras14@gmail.com"

EXPOSE 8000

RUN apk add build-base

WORKDIR /go/src/github.com/youoffcrawler

COPY go.mod go.sum ./

RUN go mod download