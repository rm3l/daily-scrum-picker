# Build stage
FROM docker.io/library/golang:1.25.0-alpine AS builder

WORKDIR /app

COPY LICENSE ./

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY pick_next.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o daily-scrum-picker pick_next.go

# Runtime stage
FROM docker.io/library/alpine:3.22

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/daily-scrum-picker .

# Create a non-root user
RUN adduser -D -s /bin/sh scrummaster
USER scrummaster

# Set the entrypoint
ENTRYPOINT ["./daily-scrum-picker"]
