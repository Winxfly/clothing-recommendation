FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /clothing-recommendation ./cmd/clothing-recommendation

# Продакшн образ
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /clothing-recommendation .
COPY config/local.yaml ./config/

EXPOSE 8082
CMD ["./clothing-recommendation"]