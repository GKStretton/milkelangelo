FROM alpine:latest
RUN apk add v4l-utils

WORKDIR /src
COPY . .

ENTRYPOINT ./configure