# Start from a base Go image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /companyservice

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o companyservice

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./companyservice"]
