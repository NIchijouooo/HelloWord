# syntax=docker/dockerfile:1
FROM golang:alpine AS builder
WORKDIR $GOPATH/src/app
ADD . ./
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
RUN go build

FROM alpine:latest
WORKDIR /go/src/app
COPY --from=builder /go/src/app .
EXPOSE 9999
RUN chmod +x goIDC
RUN pwd
RUN ls -lrt
ENTRYPOINT [ "./goIDC" ]
