# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# ---- Runtime stage ----
FROM alpine:latest

# CA certs needed for HTTPS calls (Google OAuth, MinIO, Midtrans)
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
