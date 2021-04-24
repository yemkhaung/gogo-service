FROM golang:1.16
RUN go version

ENV APP_NAME="gogo-service"

COPY . /go/src/${APP_NAME}
WORKDIR /go/src/${APP_NAME}

RUN go install .

ENTRYPOINT ${APP_NAME}
