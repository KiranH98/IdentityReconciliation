# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /IdentityReconciliation

# Copy the local source files to the container's working directory
COPY . .

# Install sqlite3 package
RUN apt-get update && apt-get install -y sqlite3

# Download and install Go dependencies
RUN go mod download

# Build the Go application
RUN go build -o main

# Expose the port your application listens on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
