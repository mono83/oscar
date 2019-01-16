FROM golang:1.11.4-alpine3.8 as builder
WORKDIR /go/src/github.com/mono83/oscar
COPY . .
RUN apk add --no-cache git make && make release

FROM alpine:3.8
COPY --from=builder /go/src/github.com/mono83/oscar/release/oscar_linux /oscar
RUN chmod a+x /oscar
RUN mkdir /data
VOLUME /data

ENTRYPOINT ["/oscar"]
