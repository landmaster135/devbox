# syntax=docker/dockerfile:1

# Build stage: compile the HTTP server binary
FROM golang:1.25-bookworm AS builder
WORKDIR /workspace

# Download dependencies first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source tree and build the Cloud Run entrypoint
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /workspace/bin/server ./cmd/http

# Runtime stage: minimal image for Cloud Run
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

# Copy the compiled binary only
COPY --from=builder /workspace/bin/server /app/server

# Cloud Run expects the service to listen on $PORT (default 8080)
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/server"]
