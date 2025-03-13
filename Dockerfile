# Use Golang base image
FROM golang:1.24

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Install dependencies
RUN go mod tidy

# Build the app
RUN go build -o main .

# Run the application
CMD ["./main"]
