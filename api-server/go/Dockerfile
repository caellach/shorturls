FROM golang:1.22-bookworm as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

# Intermediate container for updating CA certificates
FROM debian:bookworm-slim as certs
RUN apt-get update && apt-get install -y ca-certificates

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app /app
# Copy the updated CA certificates from the intermediate container
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["./app"]