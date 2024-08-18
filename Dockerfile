FROM golang:1.22

WORKDIR /usr/src/app

COPY . .

RUN go get
RUN go build

CMD ["./shortify"]
