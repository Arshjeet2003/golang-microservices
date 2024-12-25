# Stage 1: Build the Go binary using a Go base image
FROM golang:1.23-alpine AS build

# Install necessary build tools
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory inside the container
WORKDIR /go/src/github.com/Arshjeet2003/golang-microservices

# Copy the go module files (go.mod, go.sum)
COPY go.mod go.sum ./

# Copy the vendor directory (if you're using vendored dependencies)
COPY vendor/ vendor/

# Copy the source code into the container
COPY account/ account/

# Download Go modules and build the binary
RUN GO111MODULE=on go mod download
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./account/cmd/account

# Stage 2: Create a minimal final image
FROM alpine:3.11

# Set the working directory for the final image
WORKDIR /usr/bin

# Copy the built binary from the build stage
COPY --from=build /go/bin/app .

# Expose the port your application will run on
EXPOSE 8080

# Define the command to run the Go binary when the container starts
CMD ["app"]
