# Use the official Golang image to build the Go binary
FROM golang:1.19 as builder

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy all files from the current directory to the working directory in the image
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Use the official AWS Lambda Go image as the base image for the final stage
FROM public.ecr.aws/lambda/go:latest

# Copy the compiled Go binary from the builder stage to the /var/task directory in the final stage
COPY --from=builder /go/src/app/main /var/task/main

# Set the CMD to run the compiled binary as the Lambda handler
CMD [ "./main" ]
