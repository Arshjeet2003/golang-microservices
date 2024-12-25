# Use a Go base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /go/src/github.com/Arshjeet2003/golang-microservices

# Copy go module files (go.mod, go.sum)
COPY go.mod go.sum ./

# Copy the vendor directory if you're using it
COPY vendor/ vendor/

# Install dependencies (if using Go modules)
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary (ensure to use the correct Go build path)
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./catalog/cmd/catalog

# Create a new stage for the final image (to keep it minimal)
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /go/bin/app /usr/local/bin/app

# Expose the necessary port (change this based on your app's config)
EXPOSE 8080

# Command to run your Go binary when the container starts
CMD ["app"]

