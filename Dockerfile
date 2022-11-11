# Start from golang base image to build the server
FROM golang:1.19-alpine as builder

# Tools needed to compile go applications
RUN apk update && apk add --no-cache git curl make gcc g++ musl-dev binutils autoconf automake libtool pkgconfig check-dev file patch

# Set the current working directory inside the container
WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
# Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . ./

# Build the Go app
RUN make build



# Start a new stage from scratch
FROM dsuite/alpine-base

# Add Maintainer info
LABEL maintainer="Jocelyn GENNESSEAUX"

# Update current image
RUN apk add --no-cache ca-certificates && update-ca-certificates

# Define working dir
WORKDIR /mailer

# Copy the Pre-built binary file from the previous stage.
# Don't forget to copy the .env file
COPY --from=builder /build/build/mailer mailer
COPY --from=builder /build/assets/ assets/

# Expose port 8080 to the outside world
EXPOSE 8080

# declare the volume to store the list of users
VOLUME [ "/mailer/data" ]

# Test the container to check that it is still working
HEALTHCHECK --interval=5m --timeout=30s --start-period=5s --retries=10 \
    CMD curl -f http://localhost:8080/check || exit 1

# Command to run the executable
CMD ["/mailer/mailer"]
