# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o recipe-social-backend main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/recipe-social-backend .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./recipe-social-backend"]
