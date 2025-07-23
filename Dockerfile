# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY pick_next.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scrum-picker pick_next.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/scrum-picker .

# Copy the default team file
COPY team.txt ./

# Create a non-root user
RUN adduser -D -s /bin/sh scrumuser
USER scrumuser

# Set the entrypoint
ENTRYPOINT ["./scrum-picker"] 