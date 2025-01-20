FROM golang:1.23.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gostockly ./cmd/api

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/gostockly .
EXPOSE 8080
CMD ["./gostockly"]
