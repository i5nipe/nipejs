# Use the latest Go base image
FROM golang:latest

# Install nipejs
RUN go install github.com/i5nipe/nipejs/v2@latest

# Add the Go bin directory to the PATH
ENV PATH=$PATH:/go/bin

