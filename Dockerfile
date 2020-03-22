FROM golang:1.14.1-alpine AS builder
RUN apk add --no-cache make git
WORKDIR /go/src/github.com/mackerelio/mkr/
COPY . .
ENV GO111MODULE=on
RUN make build

FROM alpine:3.11.3
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/mackerelio/mkr/mkr /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/mkr"]
