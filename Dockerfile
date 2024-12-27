# Use an official Go image
FROM golang:1.23.3-alpine3.20


# Set the working directory
WORKDIR /app

COPY . .

RUN go mod tidy

# Build the application
RUN go build -o websocket-chat .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./websocket-chat"]