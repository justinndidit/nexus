# STAGE 1: Build the binary
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Leverage Docker cache for dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a statically linked binary (no CGO needed for our pgx driver)
RUN CGO_ENABLED=0 GOOS=linux go build -o /nexus-ledger ./cmd/main.go

# STAGE 2: Create the production image
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /nexus-ledger /nexus-ledger

# Use a non-root user for security
USER nonroot:nonroot

ENTRYPOINT ["/nexus-ledger"]