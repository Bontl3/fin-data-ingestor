# Build stage
FROM golang:1.17 AS builder

WORKDIR /data_ingestion_microservice
COPY . /data_ingestion_microservice/

RUN CGO_ENABLED=0 go build -o myingestorapp ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /data_ingestion_microservice/myingestorapp .
EXPOSE 8080

CMD ["./myingestorapp"]



