FROM golang:1.23-alpine AS builder

# Явно указываем прокси и отключаем IPv6
# ARG GOPROXY=https://goproxy.cn,direct
# ENV GOSUMDB=off

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]