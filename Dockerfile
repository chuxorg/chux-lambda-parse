# Use the official Golang image as the base image
FROM golang:1.19 AS builder

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o parseHandler main.go

# Use a minimal image for the final container
FROM public.ecr.aws/lambda/provided:al2

# Copy the built binary from the builder stage
COPY --from=builder /app/parseHandler /var/task/parseHandler
COPY --from=builder /app/bootstrap /var/task/bootstrap

# Set the CMD to your handler
CMD [ "parseHandler" ]
