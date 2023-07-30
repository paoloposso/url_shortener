# Use the official Go image as the base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the Go API code from the host into the container
COPY . .

# Download Go module dependencies
RUN go get -d -v ./...

# Build the Go API binary
RUN go build -o app

# Expose the port that the Go API will listen on
EXPOSE 8080

# Command to run the Go API binary
CMD ["./app"]
