FROM golang:1.20-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Copy local code to the container image.
COPY . ./

# Retrieve application dependencies.
COPY server/go.* ./
RUN go mod download

# Build the binary.
RUN go build -mod=readonly -v -o runserver server/cmd/main.go

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/runserver /app/runserver

# Run the web service on container startup.
CMD ["/app/runserver"]
