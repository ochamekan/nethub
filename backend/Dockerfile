FROM golang:1.26.1

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENV PATH="/go/bin:${PATH}"

ENTRYPOINT goose up && app

