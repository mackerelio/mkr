FROM golang:1.16-alpine AS builder
RUN apk add --no-cache make git
WORKDIR /go/src/github.com/mackerelio/mkr/
COPY . .
RUN make build

FROM alpine:3.14.1
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/mackerelio/mkr/mkr /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/mkr"]
