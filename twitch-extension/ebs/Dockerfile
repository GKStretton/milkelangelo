# Build stage
FROM golang:1.21 AS builder

WORKDIR /src
COPY ./ebs/go.mod ./go.mod
COPY ./ebs/go.sum ./go.sum
RUN go mod download
COPY ./ebs .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -a -ldflags '-extldflags "-static"'


# Runtime stage
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y
WORKDIR /app

COPY --from=builder /src/ebs .
COPY --from=builder /src/run.sh .

CMD ["./run.sh"]