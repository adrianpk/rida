# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o beak main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/beak ./
EXPOSE 8080
CMD ["./beak"]
