FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy the entire workspace if using go.work, or just the necessary files
# For simplicity in this example, we assume we are building from the root of the repo
COPY . .

WORKDIR /app/examples/docker-compose
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/examples/docker-compose/main .
COPY --from=builder /app/examples/docker-compose/config ./config

CMD ["./main"]
