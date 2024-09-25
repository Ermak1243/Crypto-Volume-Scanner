# Use the official Golang image based on Alpine Linux as the builder stage
FROM golang:alpine as builder

# Set the working directory inside the container
WORKDIR /src/app

# Copy the Go module files to the working directory
COPY ./go.mod ./go.sum . ./

# Download the Go module dependencies specified in go.mod
RUN go mod download

# Build the Go application located in cmd/app and output the binary as 'main'
RUN cd cmd/app && go build -o main

# Start a new stage from the latest Alpine image for a lightweight final image
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /

# Copy the built binary from the builder stage to the current working directory
COPY --from=builder /src/app ./cvs

# Update the package list and install Go, godoc, and make in the final image
RUN apk update && \
    apk add go && \
    go install -v golang.org/x/tools/cmd/godoc@latest && \
    apk add make

# Copy the godoc binary from its installation path to a bin directory
RUN ["cp", "./root/go/bin/godoc", "./bin"]

# Set the entry point for the container to use 'make'
ENTRYPOINT ["make"]

# Provide default arguments for 'make', specifying to run in the 'liquid' directory with parallel jobs
CMD [ "-C", "./cvs", "-j", "parallel" ]