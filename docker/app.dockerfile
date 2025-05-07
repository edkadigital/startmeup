FROM golang:1.24 AS builder

ARG VERSION=${VERSION}

WORKDIR /go/src/app

COPY . .

RUN go generate ./ent/schema

RUN CGO_ENABLED=0 go build -o startmeup -ldflags=-X=main.version=${VERSION} cmd/web/main.go

FROM alpine

RUN apk add ca-certificates

COPY --from=builder /go/src/app/startmeup /startmeup

COPY config/config.yaml /config/config.yaml

COPY static/ /static

CMD ["/startmeup"]
