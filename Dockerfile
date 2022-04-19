FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build ./cmd/simple-multiplexer/simple_multiplexer.go

EXPOSE 8080


