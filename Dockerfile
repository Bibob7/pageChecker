FROM golang:alpine3.12 as main

RUN apk add --update --no-cache \
    git \
    ca-certificates \
    jq

COPY ./src /application

WORKDIR /application

# Install the Go app
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo .
