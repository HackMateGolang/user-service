FROM golang:1.25.5-alpine AS builder
WORKDIR /build
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o user-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from-builder /build/user-service .
CMD ["./user-service"]