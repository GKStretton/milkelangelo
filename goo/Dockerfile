FROM golang:1.21

RUN apt update && apt install -y \
    ffmpeg \
    gphoto2

WORKDIR /src

COPY ./goo/go.mod ./go.mod
COPY ./goo/go.sum ./go.sum
COPY ./asol-protos ../asol-protos
COPY ./resources ../resources

RUN go mod download

COPY ./goo .

RUN go build -race

ENTRYPOINT ["./goo"]