ARG GO_VERSION=1.15.5

FROM golang:${GO_VERSION}-alpine AS builder

RUN mkdir -p /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN apk --no-cache add tzdata
RUN go build -o ./app ./main.go

FROM alpine:latest

RUN mkdir -p /api
WORKDIR /api

COPY --from=builder /api/app .
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asial/Seoul

EXPOSE 3000
ENTRYPOINT ["./app"]

