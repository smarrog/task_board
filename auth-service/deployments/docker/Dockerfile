FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go build -o auth-service ./cmd/auth-service

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/auth-service /app/auth-service

ENTRYPOINT ["/app/auth-service"]