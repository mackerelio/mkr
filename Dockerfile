FROM golang:alpine3.7 AS builder
RUN apk add --no-cache make git ca-certificates
WORKDIR /go/src/github.com/mackerelio/mkr/
COPY . .
RUN make build

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/mackerelio/mkr/mkr /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/mkr"]
