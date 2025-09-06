FROM golang:1.25-trixie AS base

WORKDIR "/build"

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o salty-spitoon ./cmd/api

EXPOSE 8080

ENTRYPOINT ["/build/salty-spitoon"]
