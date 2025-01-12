FROM golang:1.23-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Expose ports
EXPOSE 8080

# Run the application
CMD ["./main"]
