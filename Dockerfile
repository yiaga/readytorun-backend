# Dockerfile

# Stage 1: Build the Go application
# Use a Go base image with your specified version.
FROM golang:1.25.0-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to cache dependencies.
COPY go.mod .
COPY go.sum .

# Download Go modules.
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the Go application for a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

# Stage 2: Create the final, minimal image
# Use a minimal base image to reduce the final image size and attack surface.
FROM alpine:3.20

# Set the working directory.
WORKDIR /

# Copy the built binary from the 'builder' stage.
COPY --from=builder /server /server

# Expose the port the server will run on.
EXPOSE 8080

# Set environment variables as placeholders. In production, use GCP's secret management.
ENV DATABASE_URL=postgres://user:password@host:port/dbname?sslmode=disable
ENV AUTH_SECRET_KEY=a_super_secret_key

# The command to run the application.
CMD ["/server"]