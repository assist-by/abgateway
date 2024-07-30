FROM golang:1.22-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 50051 to the outside world
EXPOSE 50051

# Command to run the executable
CMD ["./main"]