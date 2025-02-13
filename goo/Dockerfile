# Build stage
FROM golang:1.21 AS builder
WORKDIR /src
COPY ./goo/go.mod ./go.mod
COPY ./goo/go.sum ./go.sum
COPY ./asol-protos ../asol-protos
COPY ./resources ../resources
RUN go mod download
COPY ./goo .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -a -ldflags '-extldflags "-static"'

# Runtime stage
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y \
    ffmpeg \
    gphoto2
WORKDIR /app
COPY --from=builder /src/goo .
COPY ./scripts ./scripts
ENTRYPOINT ["./goo"]