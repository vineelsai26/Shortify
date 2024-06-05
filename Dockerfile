FROM golang:1.22 AS builder

WORKDIR /tmp/build

COPY . .

RUN go get
RUN go build

CMD ["./shortify"]
