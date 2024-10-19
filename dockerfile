# Use the official golang image as a base
FROM golang:1.17 AS builder

# Set the working directory
WORKDIR /workspace

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the go source
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o admission-controller main.go

# Use a minimal base image
FROM alpine:3.14
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /workspace/admission-controller .

# Command to run the admission controller
ENTRYPOINT ["./admission-controller"]
