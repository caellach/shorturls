FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/app /app
ENTRYPOINT ["./app"]