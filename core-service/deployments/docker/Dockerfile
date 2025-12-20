FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go build -o core-service ./cmd/core-service

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/core-service /app/core-service

ENTRYPOINT ["/app/core-service"]