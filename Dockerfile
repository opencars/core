FROM golang:1.19-alpine AS build

WORKDIR /go/src/app

LABEL maintainer="ashanaakh@gmail.com"

RUN apk add bash ca-certificates git gcc g++ libc-dev curl make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN export BLDDIR=/go/bin && \
    make clean && \
    make build

FROM alpine

RUN apk update && apk upgrade

WORKDIR /app

COPY --from=build /go/bin/ ./
