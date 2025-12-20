FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go build -o api-service ./cmd/api-service

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/api-service /app/api-service

ENTRYPOINT ["/app/api-service"]