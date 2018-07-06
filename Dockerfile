FROM golang:alpine3.7
RUN apk add --no-cache make git
WORKDIR /go/src/github.com/mackerelio/mkr/
COPY . .
RUN make build

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/mackerelio/mkr/mkr /usr/local/bin/
ENTRYPOINT ["mkr"]
