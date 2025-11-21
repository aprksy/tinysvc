# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -o tinysvc cmd/server/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

# Copy binary and web files
COPY --from=builder /app/tinysvc .
COPY --from=builder /app/web ./web

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Environment variables
ENV SERVER_PORT=8080
ENV DB_PATH=/app/data/tinysvc.db

CMD ["./tinysvc"]