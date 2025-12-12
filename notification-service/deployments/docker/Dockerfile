FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go build -o notification-service ./cmd/notification-service

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/notification-service /app/notification-service

ENTRYPOINT ["/app/notification-service"]