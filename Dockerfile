# syntax=docker/dockerfile:1
# ------------------------
# Builder stage
FROM golang:1.25-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go.mod and go.sum for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build Go app (statically linked)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# ------------------------
# Final image
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy built binary
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the app
ENTRYPOINT ["./main"]