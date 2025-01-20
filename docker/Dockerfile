# Build stage
FROM golang:1.23.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gostockly ./cmd/api

# Runtime stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/gostockly ./
RUN apk add --no-cache ca-certificates
EXPOSE 8080
CMD ["./gostockly"]
