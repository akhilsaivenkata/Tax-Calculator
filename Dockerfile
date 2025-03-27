# --- Stage 1: Build the Go binary ---
    FROM golang:1.23.4 AS builder

    WORKDIR /app
    
    # Copy go mod and sum files
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the source code
    COPY . .
    
    # Build the Go app
    RUN CGO_ENABLED=0 GOOS=linux go build -o tax-server ./cmd/server
    
    # --- Stage 2: Slim runtime image ---
    FROM alpine:latest
    
    # Install minimal CA certificates
    RUN apk --no-cache add ca-certificates
    
    WORKDIR /app
    
    # Copy only the binary from builder
    COPY --from=builder /app/tax-server .
    
    # Expose the API port
    EXPOSE 8080
    
    # Run the Go binary
    CMD ["./tax-server"]
    