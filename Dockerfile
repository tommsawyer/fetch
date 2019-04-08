FROM golang:1.12.1 as builder
WORKDIR /go/src/github.com/tommsawyer/fetch
COPY . .
ENV GO111MODULE=on
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 make build

FROM alpine:latest
MAINTAINER Denis Maximov <mortraineymusic@gmail.com>
COPY --from=builder /go/src/github.com/tommsawyer/fetch/build/fetch /usr/bin/fetch
EXPOSE 9090
ENTRYPOINT "fetch"
