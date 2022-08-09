FROM golang:1.18.1 as builder
RUN apt-get update
RUN apt-get install unzip

WORKDIR /shakes

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Build shakes-server.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o shakes-server ./cmd/shakes-server

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /shakes/shakes-server /shakes-server

# Run the service on container startup.
CMD ["/shakes-server"]
