FROM golang:alpine

WORKDIR /go/go.adphi.net/go-repo

COPY go.mod .

RUN go mod download

