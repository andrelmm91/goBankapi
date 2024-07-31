# Build stage
FROM golang:1.22.4-alpine3.20 AS builder
WORKDIR /app

# Copy all files, including Swagger YAML or JSON and source code
COPY . .

# Build the Go application
RUN go build -o main main.go

# Run stage
FROM alpine:3.20
WORKDIR /app

# Copy the built Go binary from the build stage
COPY --from=builder /app/main .

# Copy environment files and start scripts
COPY app.env .
COPY app_dev.env .
COPY start.sh .
# COPY wait-for.sh .

# Copy the Swagger JSON file into the final image
COPY --from=builder /app/docs/swagger/swagger.json /app/swagger/swagger.json

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/app/main"]
