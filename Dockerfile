# Build GPtt in a stock Go builder container
FROM golang:1.11.1-alpine3.8 as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

RUN mkdir -p /go/src/github.com/ailabstw
ADD . /go/src/github.com/ailabstw/go-pttai
RUN cd /go/src/github.com/ailabstw/go-pttai && make gptt

# Pull GPtt into a second stage deploy alpine container
FROM alpine:3.8

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/ailabstw/go-pttai/build/bin/gptt /usr/local/bin/

EXPOSE 9774 14779 9487 9487/udp
ENTRYPOINT ["gptt", "--rpcaddr", "0.0.0.0"]
