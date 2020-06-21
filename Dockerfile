FROM golang:1.14-alpine

LABEL mantainer="Auras Popescu popescuauras14@gmail.com"

EXPOSE 8000

WORKDIR /go/src/github.com/youoffcrawler

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./main.go

CMD ["./main"]