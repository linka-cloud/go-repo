FROM golang:alpine as builder

WORKDIR /go/go.adphi.net/go-repo

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o go-repo .

FROM alpine

RUN apk add ca-certificates

COPY --from=builder /go/go.adphi.net/go-repo/go-repo /usr/bin/

USER nobody

EXPOSE 8888

ENTRYPOINT ["/usr/bin/go-repo"]
