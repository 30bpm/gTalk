ARG GO_VERSION=1.15

FROM golang:${GO_VERSION}-alpine AS builder

RUN mkdir -p /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./app ./main.go

FROM alpine:latest

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/app .

EXPOSE 3000

ENTRYPOINT ["./app"]