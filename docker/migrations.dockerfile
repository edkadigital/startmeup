FROM golang:1.24-alpine3.21

# Install Atlas CLI
RUN apk add --no-cache ca-certificates curl make
RUN curl -L https://release.ariga.io/atlas/atlas-linux-amd64-latest -o atlas
RUN chmod +x atlas && mv atlas /usr/local/bin/atlas

WORKDIR /go/src/app

COPY . .

RUN go generate ./ent/schema

RUN CGO_ENABLED=0 go build -o bin/migrate -ldflags=-X=main.version=${VERSION} cmd/migrate/main.go

RUN mv bin/migrate /usr/local/bin/migrate
