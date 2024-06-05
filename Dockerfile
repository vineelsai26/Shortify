FROM golang:1.22 AS builder

WORKDIR /tmp/build

COPY . .

RUN go get
RUN go build

FROM vineelsai/alpine AS runner

WORKDIR /usr/src/app

COPY --from=builder /tmp/build/shortify .
COPY --from=builder /tmp/build/static static

CMD ["./shortify"]
