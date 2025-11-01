# ==============================
# Stage 1: Build the Go binary
# ==============================
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build binary
RUN go build -o main .

# ==============================
# Stage 2: Run the application
# ==============================
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run
CMD ["/app/main"]
