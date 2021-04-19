FROM golang:1.16
RUN go version

ENV SERVICE_NAME="gogo-service"

COPY . /go/src/gogo-service
WORKDIR /go/src/gogo-service

RUN go install ./...

ENTRYPOINT $SERVICE_NAME
