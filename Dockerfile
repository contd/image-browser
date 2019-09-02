# Builder image
FROM golang AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server

# Production image
FROM golang:1.12-stretch
COPY --from=builder /app/server /app/

EXPOSE 6969
ENTRYPOINT [ "/app/server" ]
