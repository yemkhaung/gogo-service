FROM golang:1.16
RUN go version

COPY . /go/src/gogo-service
WORKDIR /go/src/gogo-service

RUN go build ./...

ENTRYPOINT ./gogo-service
