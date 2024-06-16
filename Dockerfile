FROM golang:1.22-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

CMD ["air", "-c", ".air.toml"]

