FROM golang:1.22-bookworm as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app /app
ENTRYPOINT ["./app"]