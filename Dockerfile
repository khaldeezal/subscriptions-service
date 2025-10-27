# Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /subscriptions ./cmd/server

# Runtime
FROM alpine:3.20
WORKDIR /
COPY --from=builder /subscriptions /subscriptions
ENV APP_PORT=8080
EXPOSE 8080
ENTRYPOINT ["/subscriptions"]
