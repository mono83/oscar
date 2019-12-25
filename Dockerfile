FROM golang:1.13.5-alpine3.10 as builder
WORKDIR /go/src/github.com/mono83/oscar
COPY . .
RUN apk add --no-cache git make && make release

FROM alpine:3.10
COPY --from=builder /go/src/github.com/mono83/oscar/release/oscar_linux /oscar
RUN chmod a+x /oscar
RUN mkdir /data
VOLUME /data

ENTRYPOINT ["/oscar"]
