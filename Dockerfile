FROM golang:buster AS builder

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o favourate main.go

EXPOSE 80

ENTRYPOINT ["/app/favourate"]