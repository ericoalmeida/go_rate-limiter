FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/

FROM alpine:3.21.3

ENV PORT=3000

WORKDIR /server

COPY --from=builder /app/app .

EXPOSE ${PORT}

CMD ["./app"]