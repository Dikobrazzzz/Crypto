FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /app/main ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]