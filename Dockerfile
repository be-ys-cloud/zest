FROM golang:1.18-bullseye

RUN apt install git

RUN mkdir /go/src/zest

WORKDIR /go/src/zest

COPY . .

RUN go mod download
RUN go build -o /go/bin/zest ./cmd

FROM debian:10
RUN mkdir -p /app

RUN apt update && apt install -y sed grep bash

COPY --from=0 /go/bin/zest /app/zest

RUN ls /app

WORKDIR /app
ENTRYPOINT ["./zest"]

EXPOSE 80
