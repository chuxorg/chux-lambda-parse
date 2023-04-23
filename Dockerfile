# # Use the official Golang image as the base image
# FROM golang:1.19 AS builder

# # Set the working directory
# WORKDIR /app

# # Copy the source code into the container
# COPY . .

# # Build the Go application
# RUN CGO_ENABLED=0 GOOS=linux go build -o parseHandler main.go

# # Use a minimal image for the final container
# FROM public.ecr.aws/lambda/provided:al2

# # Copy the built binary from the builder stage
# COPY --from=builder /app/parseHandler /var/task/parseHandler
# COPY --from=builder /app/bootstrap /var/task/bootstrap

# # Set the CMD to your handler
# CMD [ "parseHandler" ]
# Use the official AWS Lambda Go image as the base image
# Use the official AWS Lambda Go image as the base image
FROM public.ecr.aws/lambda/go:latest

# Set the working directory to /var/task
WORKDIR /var/task

# Copy all files from the current directory to the working directory in the image
COPY . .

# Build the Go binary
RUN go build -o main main.go

# Set the CMD to run the compiled binary as the Lambda handler
CMD [ "./main" ]
