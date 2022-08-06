# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This Dockerfile builds an image that runs the registry-server with no proxy.

# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.18.1 as builder
RUN apt-get update
RUN apt-get install unzip

WORKDIR /

# Get the latest registry source code
RUN git clone https://github.com/timburks/shakes

# Work in the shakes directory
WORKDIR /shakes

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
RUN go mod download

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
