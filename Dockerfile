# Use a Golang base image
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the application code into the container
COPY . .

# Build the application inside the container
RUN go build -o server .

# Create a new image
FROM alpine:latest

WORKDIR /app

# Copy the application executable from the previous stage
COPY --from=builder /app/server .

# Copy the environment variables file
COPY .env .

# Set the entrypoint for the container
ENTRYPOINT ["./server"]